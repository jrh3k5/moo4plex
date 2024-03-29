package media

import (
	"context"
	"fmt"
	"strings"
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

func (g *GORMTagService) AddTagToMetadataItem(ctx context.Context, metadataItemID int64, tagID int64, tagType gormmodel.TagType) error {
	maxIndexSelectQuery := `SELECT IFNULL(MAX("index"), 0)
						    FROM taggings
						    INNER JOIN tags ON tags.id = taggings.tag_id AND tags.tag_type = ?
						    WHERE taggings.metadata_item_id = ?`
	var maxIndex int64
	if getMaxIndexErr := g.db.WithContext(ctx).Raw(maxIndexSelectQuery, int(tagType), metadataItemID).Scan(&maxIndex).Error; getMaxIndexErr != nil {
		return fmt.Errorf("failed to get max index for metadata item ID %d, tag type %v: %w", metadataItemID, tagType, getMaxIndexErr)
	}

	insertQuery := `INSERT INTO taggings(metadata_item_id, tag_id, "index", created_at) VALUES(?, ?, ?, ?)`
	if insertErr := g.db.WithContext(ctx).Exec(insertQuery, metadataItemID, tagID, maxIndex+1, time.Now()).Error; insertErr != nil {
		return fmt.Errorf("failed to insert tag association of tag ID %d to metadata item ID %d at index %d: %w", tagID, metadataItemID, maxIndex+1, insertErr)
	}
	return nil
}

// DeleteTagAssociations deletes all tag associations for the given tag
func (g *GORMTagService) DeleteTagAssociations(ctx context.Context, tagID int64, tagType gormmodel.TagType) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		metadataIDs, metadataSelectErr := getTagMetadataIDs(ctx, tx, []int64{tagID}, tagType)
		if metadataSelectErr != nil {
			return fmt.Errorf("failed to select metadata IDs for tag ID %d, type %d: %w", tagID, int(tagType), metadataSelectErr)
		}

		if deleteErr := deleteTagAssociations(ctx, tx, metadataIDs, []int64{tagID}); deleteErr != nil {
			return fmt.Errorf("failed to delete association of tag ID %d to %d metadata items: %w", tagID, len(metadataIDs), deleteErr)
		}

		return nil
	})
}

// GetMetadataItemsForTags gets the metadata items associated to the given tags
func (g *GORMTagService) GetMetadataItemsForTags(ctx context.Context, tagIDs []int64) ([]*gormmodel.MetadataItem, error) {
	var metadataItems []*gormmodel.MetadataItem
	if dbErr := g.db.WithContext(ctx).Select("metadata_items.id, metadata_items.title, metadata_items.library_section_id").
		Joins("INNER JOIN taggings on taggings.metadata_item_id = metadata_items.id AND taggings.tag_id IN (?)", tagIDs).
		Find(&metadataItems).Error; dbErr != nil {
		return nil, fmt.Errorf("failed to look up metadata items for %d tags: %w", len(tagIDs), dbErr)
	}
	return metadataItems, nil
}

// GetMetadataItemsForTagSubstring gets all metadata items that match the given substring in their tags of the given types.
func (g *GORMTagService) GetMetadataItemsForTagSubstring(ctx context.Context, librarySectionID int64, tagTypes []gormmodel.TagType, tagTextSubstring string) ([]*gormmodel.MetadataItem, error) {
	tagTypeInts := make([]int64, len(tagTypes))
	for tagTypeIndex, tagType := range tagTypes {
		tagTypeInts[tagTypeIndex] = int64(tagType)
	}

	var metadataItems []*gormmodel.MetadataItem
	if dbErr := g.db.WithContext(ctx).Distinct("metadata_items.id, metadata_items.title, metadata_items.library_section_id").
		Joins("INNER JOIN taggings on taggings.metadata_item_id = metadata_items.id").
		Joins("INNER JOIN tags ON tags.id = taggings.tag_id AND tags.tag_type IN (?) and lower(tags.tag) LIKE ?", tagTypeInts, "%"+strings.ToLower(tagTextSubstring)+"%").
		Where("metadata_items.library_section_id = ?", librarySectionID).
		Find(&metadataItems).Error; dbErr != nil {
		return nil, fmt.Errorf("failed to look up metadata items for %d tag types with a substring of '%s': %w", len(tagTypes), tagTextSubstring, dbErr)
	}
	return metadataItems, nil
}

// GetTagsForLibrarySection gets the tags for the given library section ID and tag type
func (g *GORMTagService) GetTagsForLibrarySection(ctx context.Context, tagType gormmodel.TagType, librarySectionID int64) ([]*gormmodel.Tag, error) {
	var tags []*gormmodel.Tag
	queryDB := g.db.WithContext(ctx).Distinct("tags.id, tags.tag, tags.tag_type, metadata_items.library_section_id, tags.user_thumb_url").
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
	queryDB := g.db.WithContext(ctx).Distinct("tags.id, tags.tag, tags.tag_type, metadata_items.library_section_id, tags.user_thumb_url").
		Joins("inner join taggings on taggings.tag_id = tags.id and taggings.metadata_item_id = ?", metadataItemID).
		Joins("inner join metadata_items on metadata_items.id = taggings.metadata_item_id").
		Find(&tags, "tags.tag_type = ?", int(tagType))
	if dbErr := queryDB.Error; dbErr != nil {
		return nil, fmt.Errorf("failed to resolve tags for metadata item %d, tag type %d: %w", metadataItemID, tagType, dbErr)
	}
	return tags, nil
}

// RemoveTagsFromItem will disassociate all of the given tags from the given metadata item
func (g *GORMTagService) RemoveTagsFromItem(ctx context.Context, metadataID int64, tagType gormmodel.TagType, tagIDs []int64) error {
	deleteTagsQuery := `DELETE 
						FROM taggings
						WHERE taggings.metadata_item_id = ? AND tag_id IN (?)`
	if dbErr := g.db.WithContext(ctx).Exec(deleteTagsQuery, metadataID, tagIDs).Error; dbErr != nil {
		return fmt.Errorf("failed to delete %d tags for metadata ID %d: %w", len(tagIDs), metadataID, dbErr)
	}

	return reorderTags(ctx, g.db, metadataID, tagType)
}

// ReplaceTags replaces all associations of the given toReplaceTagIDs in the given media library section with the given replacementTagID
func (g *GORMTagService) ReplaceTags(ctx context.Context, librarySectionID int64, tagType gormmodel.TagType, toReplaceTagIDs []int64, replacementTagIDs []int64) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		metadataIDs, metadataSelectErr := getTagMetadataIDs(ctx, tx, toReplaceTagIDs, tagType)
		if metadataSelectErr != nil {
			return fmt.Errorf("failed to select metadata IDs for %d replacement tags of type %d: %w", len(toReplaceTagIDs), int(tagType), metadataSelectErr)
		}

		for _, metadataID := range metadataIDs {
			// Delete the association to the tags to be merged
			if deleteErr := deleteTagAssociations(ctx, tx, []int64{metadataID}, toReplaceTagIDs); deleteErr != nil {
				return fmt.Errorf("failed to delete %d tag associations of type %d for metadata ID %d: %w", len(toReplaceTagIDs), int(tagType), metadataID, deleteErr)
			}

			// Rebuild the indices of the remaining tags to fill in gaps
			if reorderErr := reorderTags(ctx, tx, metadataID, tagType); reorderErr != nil {
				return fmt.Errorf("failed to reorder tags of type '%v' for metadata ID %d: %w", tagType, metadataID, reorderErr)
			}

			for _, replacementTagID := range replacementTagIDs {
				hasMergeTarget := false
				tagIDs, getTagIDsErr := getTagIDs(ctx, tx, metadataID, tagType)
				if getTagIDsErr != nil {
					return fmt.Errorf("failed to retrieve tag IDs of type '%v' for metadata ID %d after reordering: %w", tagType, metadataID, getTagIDsErr)
				} else {
					for _, tagID := range tagIDs {
						hasMergeTarget = tagID == replacementTagID
						if hasMergeTarget {
							break
						}
					}
				}

				// Now create the new tag, if it doesn't exist
				if !hasMergeTarget {
					tagging := &gormmodel.Tagging{
						MetadataItemID: metadataID,
						TagID:          replacementTagID,
						Index:          int64(len(tagIDs)),
						CreatedAt:      time.Now(),
					}
					if createErr := tx.Create(tagging).Error; createErr != nil {
						return fmt.Errorf("failed to add target tag to metadata ID %d: %w", metadataID, createErr)
					}
				}
			}
		}

		// Don't actually delete the tag because it has a collating tokenizer not recognized by this driver

		return nil
	})
}

// deleteTagAssociation deletes the associations between the given metadata item and tags
func deleteTagAssociations(ctx context.Context, db *gorm.DB, metadataItemIDs []int64, tagIDs []int64) error {
	deleteTaggingsQuery := `DELETE FROM taggings WHERE metadata_item_id IN (?) AND tag_id IN (?)`
	if deleteErr := db.Exec(deleteTaggingsQuery, metadataItemIDs, tagIDs).Error; deleteErr != nil {
		return fmt.Errorf("failed to delete %d tag associations for %d metadata items: %w", len(tagIDs), len(metadataItemIDs), deleteErr)
	}
	return nil
}

// getTagMetadataIDs gets the metadata item IDs that are associated to the given tag IDs
func getTagMetadataIDs(ctx context.Context, db *gorm.DB, tagIDs []int64, tagType gormmodel.TagType) ([]int64, error) {
	metadataIDSelectQuery := `SELECT DISTINCT t1.metadata_item_id
							  FROM taggings t1
							  INNER JOIN tags t2 ON t2.id = t1.tag_id AND t2.tag_type = ?
							  WHERE t1.tag_id IN (?)`
	var metadataIDs []int64
	if metadataSelectErr := db.Raw(metadataIDSelectQuery, int(tagType), tagIDs).Find(&metadataIDs).Error; metadataSelectErr != nil {
		return nil, fmt.Errorf("failed to select metadata IDs for %d replacement tags of type %d: %w", len(tagIDs), int(tagType), metadataSelectErr)
	}
	return metadataIDs, nil
}

// reorderTags pulls the tags of the given type for the given media item ID and rebuilds their indices
// so that any deleted tags' gaps are now filled and sequential ordering of the indices is maintained
func reorderTags(ctx context.Context, db *gorm.DB, metadataID int64, tagType gormmodel.TagType) error {
	tagIndexUpdateQuery := `UPDATE taggings SET "index" = ? WHERE id = ?`

	tagIDs, getTagsErr := getTagIDs(ctx, db, metadataID, tagType)
	if getTagsErr != nil {
		return fmt.Errorf("failed to get all tags for metadata ID %d after deletion: %w", metadataID, getTagsErr)
	}

	for tagIndex, tagID := range tagIDs {
		if updateErr := db.WithContext(ctx).Exec(tagIndexUpdateQuery, tagIndex, tagID).Error; updateErr != nil {
			return fmt.Errorf("failed to update tag ID %d to index %d for metadata ID %d: %w", tagID, tagIndex, metadataID, updateErr)
		}
	}

	return nil
}

func getTagIDs(ctx context.Context, db *gorm.DB, metadataID int64, tagType gormmodel.TagType) ([]int64, error) {
	getTagIDsQuery := `SELECT taggings.id FROM taggings 
					   INNER JOIN tags ON tags.id = taggings.tag_id AND tags.tag_type = ?
					   WHERE taggings.metadata_item_id = ?
					   ORDER BY "taggings.index" ASC`

	var tagIDs []int64
	if getTagsErr := db.WithContext(ctx).Raw(getTagIDsQuery, int(tagType), metadataID).Scan(&tagIDs).Error; getTagsErr != nil {
		return nil, fmt.Errorf("failed to get all tag IDs of type %v for metadata ID %d: %w", metadataID, tagType, getTagsErr)
	}
	return tagIDs, nil
}
