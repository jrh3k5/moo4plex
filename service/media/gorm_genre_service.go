package media

import (
	"context"
	"fmt"
	"time"

	"github.com/jrh3k5/moo4plex/model"
	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

const (
	genreTagType = 1
)

type GORMGenreService struct {
	db *gorm.DB
}

func NewGORMGenreService(db *gorm.DB) *GORMGenreService {
	return &GORMGenreService{
		db: db,
	}
}

func (g *GORMGenreService) GetGenres(ctx context.Context, mediaLibraryID int64) ([]*model.Genre, error) {
	var tags []*gormmodel.Tag
	queryDB := g.db.WithContext(ctx).Distinct("tags.id, tags.tag, tags.tag_type, metadata_items.library_section_id").
		Joins("inner join taggings on taggings.tag_id = tags.id").
		Joins("inner join metadata_items on metadata_items.id = taggings.metadata_item_id and metadata_items.library_section_id = ?", mediaLibraryID).
		Find(&tags, "tag_type = ?", genreTagType)
	if dbErr := queryDB.Error; dbErr != nil {
		return nil, fmt.Errorf("failed to resolve genres for media library %d: %w", mediaLibraryID, dbErr)
	}

	genres := make([]*model.Genre, len(tags))
	for tagIndex, tag := range tags {
		genres[tagIndex] = model.NewGenre(tag.ID, tag.Tag, tag.LibrarySectionID)
	}
	return genres, nil
}

func (g *GORMGenreService) MergeGenres(ctx context.Context, mergeTarget *model.Genre, toMerge []*model.Genre, totalCountCallback func(int), itemCompletionCallback func()) error {
	metadataIDSelectQuery := `SELECT DISTINCT t1.metadata_item_id
							  FROM taggings t1
							  WHERE t1.tag_id IN (?)`
	toMergeIDs := make([]int64, len(toMerge))
	for mergeIndex, mergeable := range toMerge {
		toMergeIDs[mergeIndex] = mergeable.ID
	}

	var metadataIDs []int64
	if metadataSelectErr := g.db.Raw(metadataIDSelectQuery, toMergeIDs).Scan(&metadataIDs).Error; metadataSelectErr != nil {
		return fmt.Errorf("failed to select metadata IDs for %d genres: %w", len(toMerge), metadataSelectErr)
	}

	deleteTagsQuery := `DELETE FROM tags WHERE id IN (?)`
	deleteTaggingsQuery := `DELETE FROM taggings WHERE metadata_item_id = ? AND tag_id IN (?)`
	getGenreQueries := `SELECT taggings.id FROM taggings 
						INNER JOIN tags ON tags.id = taggings.tag_id AND tags.tag_type = ?
						WHERE taggings.metadata_item_id = ?
						ORDER BY "taggings.index" ASC`
	tagIndexUpdateQuery := `UPDATE taggings SET index = ? WHERE id = ?`
	totalCountCallback(len(metadataIDs))
	for _, metadataID := range metadataIDs {
		// Delete the association to the genres to be merged
		if deleteErr := g.db.Raw(deleteTaggingsQuery, metadataID, toMergeIDs).Error; deleteErr != nil {
			return fmt.Errorf("failed to delete %d genre associations for metadata ID %d: %w", len(toMergeIDs), metadataID, deleteErr)
		}

		// Rebuild the indices of the remaining genres to fill in gaps
		var tagIDs []int64
		if getGenresErr := g.db.Raw(getGenreQueries, genreTagType, metadataID).Scan(&tagIDs).Error; getGenresErr != nil {
			return fmt.Errorf("failed to get all genres for metadata ID %d after deletion: %w", metadataID, getGenresErr)
		}

		hasMergeTarget := false
		for tagIndex, tagID := range tagIDs {
			if updateErr := g.db.Raw(tagIndexUpdateQuery, tagIndex, tagID).Error; updateErr != nil {
				return fmt.Errorf("failed to update tag ID %d to index %d for metadata ID %d: %w", tagID, tagIndex, metadataID, updateErr)
			}
			hasMergeTarget = hasMergeTarget || tagID == mergeTarget.ID
		}

		// Now create the new tag, if it doesn't exist
		if !hasMergeTarget {
			tagging := &gormmodel.Tagging{
				MetadataItemID: metadataID,
				TagID:          mergeTarget.ID,
				Index:          int64(len(tagIDs)),
				CreatedAt:      time.Now(),
			}
			if createErr := g.db.Create(tagging).Error; createErr != nil {
				return fmt.Errorf("failed to add target tag to metadata ID %d: %w", metadataID, createErr)
			}
		}
		itemCompletionCallback()
	}

	if deleteTagsErr := g.db.Raw(deleteTagsQuery, toMergeIDs).Error; deleteTagsErr != nil {
		return fmt.Errorf("failed to delete remaining %d genres: %w", len(toMergeIDs), deleteTagsErr)
	}

	return nil
}
