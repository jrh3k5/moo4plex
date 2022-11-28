package media

import (
	"context"
	"fmt"
	"time"

	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

// GORMTagService provides a reusable means of interacting with tag data in Plex
type GORMTagService struct {
	db *gorm.DB
}

func NewGORMTagService(db *gorm.DB) *GORMTagService {
	return &GORMTagService{
		db: db,
	}
}

// GetTagsForLibrarySection gets the tags for the given library section ID and tag type
func (g *GORMTagService) GetTagsForLibrarySection(ctx context.Context, tagType gormmodel.TagType, librarySectionID int64) ([]*gormmodel.Tag, error) {
	var tags []*gormmodel.Tag
	queryDB := g.db.WithContext(ctx).Distinct("tags.id, tags.tag, tags.tag_type, metadata_items.library_section_id").
		Joins("inner join taggings on taggings.tag_id = tags.id").
		Joins("inner join metadata_items on metadata_items.id = taggings.metadata_item_id and metadata_items.library_section_id = ?", librarySectionID).
		Find(&tags, "tag_type = ?", int(tagType))
	if dbErr := queryDB.Error; dbErr != nil {
		return nil, fmt.Errorf("failed to resolve tags for library section %d, tag type %d: %w", librarySectionID, tagType, dbErr)
	}
	return tags, nil
}

// GetTagsForMetadataItem gets tags of the given type for the given metadata item
func (g *GORMTagService) GetTagsForMetadataItem(ctx context.Context, tagType gormmodel.TagType, metadataItemID int64) ([]*gormmodel.Tag, error) {
	var tags []*gormmodel.Tag
	queryDB := g.db.WithContext(ctx).Distinct("tags.id, tags.tag, tags.tag_type, metadata_items.library_section_id").
		Joins("inner join taggings on taggings.tag_id = tags.id and taggings.metadata_item_id = ?", metadataItemID).
		Joins("inner join metadata_items on metadata_items.id = taggings.metadata_item_id").
		Find(&tags, "tags.tag_type = ?", int(tagType))
	if dbErr := queryDB.Error; dbErr != nil {
		return nil, fmt.Errorf("failed to resolve tags for metadata item %d, tag type %d: %w", metadataItemID, tagType, dbErr)
	}
	return tags, nil
}

// ReplaceTags replaces all associations of the given toReplaceTagIDs in the given media library section with the given replacementTagID
func (g *GORMTagService) ReplaceTags(ctx context.Context, librarySectionID int64, tagType gormmodel.TagType, toReplaceTagIDs []int64, replacementTagID int64) error {
	metadataIDSelectQuery := `SELECT DISTINCT t1.metadata_item_id
							  FROM taggings t1
							  INNER JOIN tags t2 ON t2.id = t1.tag_id AND t2.tag_type = ?
							  WHERE t1.tag_id IN (?)`

	var metadataIDs []int64
	if metadataSelectErr := g.db.Raw(metadataIDSelectQuery, int(tagType), toReplaceTagIDs).Find(&metadataIDs).Error; metadataSelectErr != nil {
		return fmt.Errorf("failed to select metadata IDs for %d replacement tags of type %d: %w", len(toReplaceTagIDs), int(tagType), metadataSelectErr)
	}

	deleteTaggingsQuery := `DELETE FROM taggings WHERE metadata_item_id = ? AND tag_id IN (?)`
	getTagsQueries := `SELECT taggings.id FROM taggings 
						INNER JOIN tags ON tags.id = taggings.tag_id AND tags.tag_type = ?
						WHERE taggings.metadata_item_id = ?
						ORDER BY "taggings.index" ASC`
	tagIndexUpdateQuery := `UPDATE taggings SET "index" = ? WHERE id = ?`
	for _, metadataID := range metadataIDs {
		// Delete the association to the tags to be merged
		if deleteErr := g.db.Exec(deleteTaggingsQuery, metadataID, toReplaceTagIDs).Error; deleteErr != nil {
			return fmt.Errorf("failed to delete %d tag associations of type %d for metadata ID %d: %w", len(toReplaceTagIDs), int(tagType), metadataID, deleteErr)
		}

		// Rebuild the indices of the remaining tags to fill in gaps
		var tagIDs []int64
		if getTagsErr := g.db.Raw(getTagsQueries, int(tagType), metadataID).Scan(&tagIDs).Error; getTagsErr != nil {
			return fmt.Errorf("failed to get all tags for metadata ID %d after deletion: %w", metadataID, getTagsErr)
		}

		hasMergeTarget := false
		for tagIndex, tagID := range tagIDs {
			if updateErr := g.db.Exec(tagIndexUpdateQuery, tagIndex, tagID).Error; updateErr != nil {
				return fmt.Errorf("failed to update tag ID %d to index %d for metadata ID %d: %w", tagID, tagIndex, metadataID, updateErr)
			}
			hasMergeTarget = hasMergeTarget || tagID == replacementTagID
		}

		// Now create the new tag, if it doesn't exist
		if !hasMergeTarget {
			tagging := &gormmodel.Tagging{
				MetadataItemID: metadataID,
				TagID:          replacementTagID,
				Index:          int64(len(tagIDs)),
				CreatedAt:      time.Now(),
			}
			if createErr := g.db.Create(tagging).Error; createErr != nil {
				return fmt.Errorf("failed to add target tag to metadata ID %d: %w", metadataID, createErr)
			}
		}
	}

	// Don't actually delete the tag because it has a collating tokenizer not recognized by this driver

	return nil
}
