package mcapi

type DatasetFileSelection struct {
	IncludeFiles []string `json:"include_files"`
	ExcludeFiles []string `json:"exclude_files"`
	IncludeDirs  []string `json:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs"`
}

type CreateOrUpdateDatasetRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Summary     string   `json:"summary"`
	License     string   `json:"license"`
	Funding     string   `json:"funding"`
	Experiments []int    `json:"experiments"`
	Communities []int    `json:"communities"`
	Tags        []Tag    `json:"tags"`
	Authors     []Author `json:"ds_authors"`
}

type Tag struct {
	Value string `json:"value"`
}

type Author struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Affiliations string `json:"affiliations"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type CreateExperimentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	ProjectID   int    `json:"project_id"`
}

type CreateActivityRequest struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	ProjectID    int         `json:"project_id"`
	ExperimentID int         `json:"experiment_id"`
	Attributes   []Attribute `json:"attributes"`
}

type CreateEntityRequest struct {
	Name         string      `json:"name"`
	Category     string      `json:"category"`
	Description  string      `json:"description"`
	Summary      string      `json:"summary"`
	ExperimentID int         `json:"experiment_id"`
	ProjectID    int         `json:"project_id"`
	ActivityID   int         `json:"activity_id"`
	Attributes   []Attribute `json:"attributes"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
	Unit  string `json:"unit"`
}

type CreateEntityStateRequest struct {
	Current    bool        `json:"current"`
	Attributes []Attribute `json:"attributes"`
}
