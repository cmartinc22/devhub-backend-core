package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pedidosya/peya-go/logs"

)

type Spec map[string]interface{}

type RelationshipType string

type ArchiveOperation bool

func (r RelationshipType) String() string {
	return string(r)
}

const (
	PUBLIC       string           = "public"
	INTERNAL     string           = "internal"
	PRIVATE      string           = "private"
	SCOPE_READ   string           = "read"
	SCOPE_WRITE  string           = "write"
	SCOPE_IGNORE string           = ""
	ARCHIVE      ArchiveOperation = true
	UNARCHIVE    ArchiveOperation = false
)

const (
	BACKEND_APP             string = "backendapp"
	MOBILE_APP              string = "mobileapp"
	WEBFRONTEND_APP         string = "webfrontend"
	DOMAIN_ENTITY           string = "domain"
	GROUP_ENTITY            string = "group"
	USER_ENTITY             string = "user"
	SERVICE_COMPONENT       string = "service"
	SCHEDULE_TASK_COMPONENT string = "scheduledtask"
	CONSUMER_COMPONENT      string = "consumer"
	MOBILE_MODULE_COMPONENT string = "mobilemodule"
	MICROFRONTEND_COMPONENT string = "microfrontend"
	API_COMPONENT           string = "api"
	DATABASE_COMPONENT      string = "database"
	TOPIC_COMPONENT         string = "topic"
	QUEUE_COMPONENT         string = "queue"
	FILE_STORAGE_COMPONENT  string = "filestorage"
	LIBRARY_STANDALONE      string = "library"
	EXECUTABLE              string = "executable"
)

const (
	ParentOf      RelationshipType = "parentOf"
	ChildOf       RelationshipType = "childOf"
	PartOf        RelationshipType = "partOf"
	HasPart       RelationshipType = "hasPart"
	DependsOn     RelationshipType = "dependsOn"
	DependencyOf  RelationshipType = "dependencyOf"
	ProvidesApi   RelationshipType = "providesApi"
	ApiProvidedBy RelationshipType = "apiProvidedBy"
	ConsumesApi   RelationshipType = "consumesApi"
	ApiConsumedBy RelationshipType = "apiConsumedBy"
	OwnedBy       RelationshipType = "ownedBy"
	OwnerOf       RelationshipType = "ownerOf"
	Unknown       RelationshipType = "UNKNOWN"
)

const (
	SOURCE         string = "source"
	TARGET         string = "target"
	TARGET_UNKNOWN string = ""
)

type Entity struct {
	Id                  string         `db:"id"`
	ApiVersion          string         `db:"api_version"`
	Kind                string         `db:"kind"`
	Slug                string         `db:"name"`
	Spec                []uint8        `db:"spec"`
	CreatedAt           time.Time      `db:"created_at"`
	ChangedAt           time.Time      `db:"changed_at"`
	Version             string         `db:"version"`
	TargetRelationships []Relationship `db:"-"`
	SourceRelationships []Relationship `db:"-"`
	LastChangedBy       string         `db:"last_changed_by"`
	CreatedBy           string         `db:"created_by"`
	Errors              []EntityError  `db:"-"`
	IsArchived          bool           `db:"is_archived"`
}

type EntityFilter struct {
	IncludeArchived *bool
	IsArchived      *bool
	ModifiedSince   *time.Time
}

type Relationship struct {
	Id             string  `db:"id"`
	Type           string  `db:"type"`
	Name           string  `db:"name"`
	Target         string  `db:"target_entity"`
	Source         string  `db:"source_entity"`
	TargetEntityID *string `db:"target_fk"`
	SourceEntityID *string `db:"source_fk"`
}

func (e *Entity) FullName() string {
	return fmt.Sprintf("%s:%s", e.Kind, e.Slug)
}

func (e *EntityWeb) FullName() string {
	return fmt.Sprintf("%s:%s", e.Kind, e.Slug)
}

func (id *EntityIdentifier) FullName() string {
	return fmt.Sprintf("%s:%s", id.Kind, id.Slug)
}

func (r *Relationship) GetTargetEntityProperties() (string, string) {
	return UnmarshalEntity(r.Target)
}

func (r *Relationship) GetSourceEntityProperties() (string, string) {
	return UnmarshalEntity(r.Source)
}

func UnmarshalEntity(entity string) (string, string) {
	parts := strings.Split(entity, ":")
	kind := parts[0]
	slug := parts[1]

	return kind, slug
}

func (e *Entity) ToWebEntity() (*EntityWeb, error) {
	relationships := []RelationshipWeb{}

	for _, relation := range e.SourceRelationships {
		webRelation, err := toWebRelationship(relation, true)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, *webRelation)
	}

	for _, relation := range e.TargetRelationships {
		webRelation, err := toWebRelationship(relation, false)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, *webRelation)
	}

	var spec map[string]interface{}

	err := json.Unmarshal(e.Spec, &spec)
	if err != nil {
		return nil, err
	}

	return &EntityWeb{
		ApiVersion: e.ApiVersion,
		Kind:       e.Kind,
		Slug:       e.Slug,
		ChangeControl: ChangeControl{
			CreatedAt:     e.CreatedAt,
			CreatedBy:     e.CreatedBy,
			ChangedAt:     e.ChangedAt,
			LastChangedBy: e.LastChangedBy,
			Version:       e.Version,
		},
		Relationships: relationships,
		Spec:          spec,
		Errors:        e.Errors,
		IsArchived:    e.IsArchived,
	}, nil
}

func toWebRelationship(dbRelationship Relationship, isSource bool) (*RelationshipWeb, error) {
	dbRelType := dbRelationship.Type // <source_type>/<target_type>, "partOf/hasPart"
	types := strings.Split(dbRelType, "/")
	if len(types) != 2 {
		return nil, fmt.Errorf("invalid relationship type")
	}

	target := ""

	relType := Unknown
	if isSource {
		relType = ParseRelationshipType(types[0])
		target = dbRelationship.Target
	} else {
		relType = ParseRelationshipType(types[1])
		target = dbRelationship.Source
	}
	if relType == Unknown {
		return nil, fmt.Errorf("Unknown relationship type")
	}

	return &RelationshipWeb{
		Name:   dbRelationship.Name,
		Target: target,
		Type:   relType,
	}, nil
}

type EntityWeb struct {
	ApiVersion    string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind          string            `json:"kind" yaml:"kind"`
	Slug          string            `json:"slug" yaml:"slug"`
	ChangeControl ChangeControl     `json:"changeControl" yaml:"changeControl"`
	Relationships []RelationshipWeb `json:"relationships,omitempty" yaml:"relationships,omitempty"`
	Spec          Spec              `json:"spec,omitempty" yaml:"spec,omitempty"`
	Errors        []EntityError     `json:"errors,omitempty" yaml:"errors,omitempty"`
	IsArchived    bool              `json:"-" yaml:"-"`
}

type ChangeControl struct {
	CreatedAt     time.Time `json:"createdAt" yaml:"createdAt"`
	ChangedAt     time.Time `json:"changedAt" yaml:"changedAt"`
	Version       string    `json:"version" yaml:"version"`
	LastChangedBy string    `json:"lastChangedBy" yaml:"lastChangedBy"`
	CreatedBy     string    `json:"createdBy" yaml:"createdBy"`
}

type RelationshipWeb struct {
	Target string           `json:"target" yaml:"target"`
	Type   RelationshipType `json:"type" yaml:"type"`
	Name   string           `json:"name" yaml:"name"`
}

func KnownKinds() []string {
	return []string{
		DOMAIN_ENTITY, GROUP_ENTITY, USER_ENTITY, BACKEND_APP, MOBILE_APP, WEBFRONTEND_APP,
		SERVICE_COMPONENT, SCHEDULE_TASK_COMPONENT, CONSUMER_COMPONENT, MOBILE_MODULE_COMPONENT, MICROFRONTEND_COMPONENT,
		API_COMPONENT, DATABASE_COMPONENT, TOPIC_COMPONENT, QUEUE_COMPONENT, FILE_STORAGE_COMPONENT, LIBRARY_STANDALONE, EXECUTABLE,
	}
}

func (e *EntityWeb) ExtractRelationships() {
	switch e.Kind {
	case DOMAIN_ENTITY:
		extractDomainRelationships(e)
	case BACKEND_APP:
		extractApplicationRelationships(e)
	case MOBILE_APP, WEBFRONTEND_APP:
		extractMobileAppRelationships(e)
	case SERVICE_COMPONENT:
		extractServiceRelationships(e)
	case SCHEDULE_TASK_COMPONENT:
		extractScheduledTaskRelationships(e)
	case CONSUMER_COMPONENT:
		extractConsumerRelationships(e)
	case MOBILE_MODULE_COMPONENT, MICROFRONTEND_COMPONENT:
		extractMicrofrontendRelationships(e)
	case API_COMPONENT, DATABASE_COMPONENT, TOPIC_COMPONENT, QUEUE_COMPONENT, FILE_STORAGE_COMPONENT:
		extractComponentRelationships(e)
	case LIBRARY_STANDALONE:
		extractLibraryRelationships(e)
	case GROUP_ENTITY, USER_ENTITY:
		logs.Debugf("[models] kind '%s' hasn't relations to extract", e.Kind)
	default:
		logs.Debugf("[models] kind '%s' not known yet", e.Kind)
	}
}

func isEntityTarget(target string) bool {
	kind, _ := UnmarshalEntity(target)
	for _, i := range KnownKinds() {
		if i == kind {
			return true
		}
	}
	return false
}

func extractRelationShips(entity *EntityWeb, rel_type RelationshipType, rel_name string) []RelationshipWeb {
	relationships := []RelationshipWeb{}
	v_interface := entity.Spec[rel_name]

	if v_interface != nil {
		switch a := v_interface.(type) {
		case []interface{}:
			for _, i := range a {
				if isEntityTarget(i.(string)) {
					relationships = append(relationships, RelationshipWeb{
						Target: i.(string),
						Type:   rel_type,
						Name:   rel_name,
					})
				}
			}
		default:
			if isEntityTarget(v_interface.(string)) {
				relationships = append(relationships, RelationshipWeb{
					Target: v_interface.(string),
					Type:   rel_type,
					Name:   rel_name,
				})
			}
		}
	}
	return relationships
}

func extractDomainRelationships(entity *EntityWeb) {
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "owner")...)
}

func extractApplicationRelationships(entity *EntityWeb) {
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, PartOf, "domain")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "owner")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "onCall")...)
}

func extractComponentRelationships(entity *EntityWeb) {
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, PartOf, "application")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "owner")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "onCall")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, DependsOn, "dependsOn")...)
}

func extractMobileAppRelationships(entity *EntityWeb) {
	extractApplicationRelationships(entity)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ConsumesApi, "consumesApis")...)
}

func extractServiceRelationships(entity *EntityWeb) {
	extractComponentRelationships(entity)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ProvidesApi, "providesApis")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ConsumesApi, "consumesApis")...)
}

func extractScheduledTaskRelationships(entity *EntityWeb) {
	extractComponentRelationships(entity)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ConsumesApi, "consumesApis")...)
}

func extractConsumerRelationships(entity *EntityWeb) {
	extractComponentRelationships(entity)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, DependsOn, "consumesFrom")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, DependsOn, "publishesTo")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ConsumesApi, "consumesApis")...)
}

func extractMicrofrontendRelationships(entity *EntityWeb) {
	extractComponentRelationships(entity)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, ConsumesApi, "consumesApis")...)
}

func extractLibraryRelationships(entity *EntityWeb) {
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, OwnedBy, "owner")...)
	entity.Relationships = append(entity.Relationships, extractRelationShips(entity, DependsOn, "dependsOn")...)
}

func (e *EntityWeb) ToDbEntity() (*Entity, error) {
	spec, err := json.Marshal(e.Spec)
	if err != nil {
		return nil, err
	}

	entity := &Entity{
		ApiVersion:    e.ApiVersion,
		Kind:          e.Kind,
		Slug:          e.Slug,
		Spec:          spec,
		CreatedAt:     e.ChangeControl.CreatedAt,
		CreatedBy:     e.ChangeControl.CreatedBy,
		ChangedAt:     e.ChangeControl.ChangedAt,
		LastChangedBy: e.ChangeControl.LastChangedBy,
		Version:       e.ChangeControl.Version,
		Errors:        e.Errors,
	}

	// Relations Section
	for _, r := range e.Relationships {
		db_relation, err := toDbRelationship(&r, e.Kind, e.Slug)
		if err != nil {
			return nil, fmt.Errorf("can't parse to DB Entity Relationships")
		}
		entity.SourceRelationships = append(entity.SourceRelationships, *db_relation)
	}
	return entity, nil
}

func toDbRelationship(relation *RelationshipWeb, kind string, slug string) (*Relationship, error) {
	target_type := TargetTypeForSources(relation.Type)

	if target_type == Unknown {
		return nil, fmt.Errorf("Unknown source relationship type")
	}

	return &Relationship{
		Type:   fmt.Sprintf("%s/%s", relation.Type, target_type),
		Name:   relation.Name,
		Target: relation.Target,
		Source: fmt.Sprintf("%s:%s", kind, slug),
	}, nil
}

func (e *EntityWeb) AddError(path string, code ErrorCode, message string) {
	e.Errors = append(e.Errors, EntityError{
		Path:    path,
		Code:    code,
		Message: message,
	})
}

func ParseRelationshipType(s string) RelationshipType {
	switch strings.ToLower(s) {
	case "parentof":
		return ParentOf
	case "childof":
		return ChildOf
	case "partof":
		return PartOf
	case "haspart":
		return HasPart
	case "dependson":
		return DependsOn
	case "dependencyof":
		return DependencyOf
	case "providesapi":
		return ApiProvidedBy
	case "apiprovidedby":
		return ProvidesApi
	case "consumesapi":
		return ConsumesApi
	case "apiconsumedby":
		return ApiConsumedBy
	case "ownedby":
		return OwnedBy
	case "ownerof":
		return OwnerOf
	default:
		return Unknown
	}
}

func TargetTypeForSources(sourceType RelationshipType) RelationshipType {
	switch sourceType {
	case OwnedBy:
		return OwnerOf
	case ChildOf:
		return ParentOf
	case PartOf:
		return HasPart
	case DependsOn:
		return DependencyOf
	case ProvidesApi:
		return ApiProvidedBy
	case ConsumesApi:
		return ApiConsumedBy
	default:
		return Unknown
	}
}

func SourceTypeForTargets(sourceType RelationshipType) RelationshipType {
	switch sourceType {
	case OwnerOf:
		return OwnedBy
	case ParentOf:
		return ChildOf
	case HasPart:
		return PartOf
	case DependencyOf:
		return DependsOn
	case ApiProvidedBy:
		return ProvidesApi
	case ApiConsumedBy:
		return ConsumesApi
	default:
		return Unknown
	}
}

func (r RelationshipType) FullRelationTypeName() string {
	switch r {
	case OwnedBy, OwnerOf:
		return fmt.Sprintf("%s/%s", OwnedBy, OwnerOf)
	case ChildOf, ParentOf:
		return fmt.Sprintf("%s/%s", ChildOf, ParentOf)
	case PartOf, HasPart:
		return fmt.Sprintf("%s/%s", PartOf, HasPart)
	case DependsOn, DependencyOf:
		return fmt.Sprintf("%s/%s", DependsOn, DependencyOf)
	case ProvidesApi, ApiProvidedBy:
		return fmt.Sprintf("%s/%s", ProvidesApi, ApiProvidedBy)
	case ConsumesApi, ApiConsumedBy:
		return fmt.Sprintf("%s/%s", ConsumesApi, ApiConsumedBy)
	default:
		return Unknown.String()
	}
}

func (r RelationshipType) Target() string {
	switch r {
	case ParentOf, PartOf, DependsOn, ApiProvidedBy, ConsumesApi, OwnedBy:
		return SOURCE
	case ChildOf, HasPart, DependencyOf, ProvidesApi, ApiConsumedBy, OwnerOf:
		return TARGET
	default:
		return TARGET_UNKNOWN
	}
}

type EntityIdentifier struct {
	Kind string
	Slug string
}

type EntitiesNews struct {
	Results  []EntityNew `json:"results" yaml:"results"`
	PageInfo Cursor      `json:"pageInfo" yaml:"pageInfo"`
}

type EntitiesSearch struct {
	Entities []EntityWeb `json:"entities" yaml:"entities"`
	PageInfo Cursor      `json:"pageInfo" yaml:"pageInfo"`
}

type EntityNew struct {
	Operation string    `json:"op" yaml:"op"`
	Entity    EntityWeb `json:"entity" yaml:"entity"`
}

type Cursor struct {
	Next *string `json:"next" yaml:"next"`
}
