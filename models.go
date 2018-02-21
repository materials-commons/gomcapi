package mcapi

type Project struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Owner       string         `json:"owner"`
	Description string         `json:"description"`
	Birthtime   Timestamp      `json:"birthtime"`
	MTime       Timestamp      `json:"mtime"`
	FileCount   int            `json:"files"`
	Notes       []*ProjectNote `json:"notes"`
	Experiments []*Experiment  `json:"experiments"`
	Samples     []*Sample      `json:"samples"`
	Todos       []*ProjectTodo `json:"todos"`
	Users       []*ProjectUser `json:"users"`
}

type ProjectNote struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Note      string    `json:"note"`
	Birthtime Timestamp `json:"-"` // `json:"birthtime"`
	MTime     Timestamp `json:"-"` // `json:"mtime"`
	Owner     string    `json:"owner"`
}

type ProjectTodo struct {
	Title string `json:"title"`
}

type ProjectUser struct {
	UserID    string    `json:"user_id"`
	Fullname  string    `json:"fullname"`
	Birthtime Timestamp `json:"-"` // `json:"birthtime"`
}

type Experiment struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Owner         string     `json:"owner"`
	Description   string     `json:"description"`
	Birthtime     Timestamp  `json:"-"` // `json:"birthtime"`
	MTime         Timestamp  `json:"-"` // `json:"mtime"`
	Citations     []string   `json:"citations"`
	Collaborators []string   `json:"collaborators"`
	Funding       []string   `json:"funding"`
	Goals         []string   `json:"goals"`
	Publications  []string   `json:"publications"`
	Papers        []string   `json:"papers"`
	Status        string     `json:"status"`
	Processes     []*Process `json:"processes"`
	Samples       []*Sample  `json:"samples"`
	Datasets      []*Dataset `json:"datasets"`
}

type Sample struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Description   string    `json:"description"`
	PropertySetID string    `json:"property_set_id"`
	Birthtime     Timestamp `json:"-"` // `json:"birthtime"`
	MTime         Timestamp `json:"-"` // `json:"mtime"`
}

type Process struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Description   string    `json:"description"`
	DoesTransform bool      `json:"does_transform"`
	ProcessType   string    `json:"process_type"`
	Birthtime     Timestamp `json:"-"` // `json:"birthtime"`
	MTime         Timestamp `json:"-"` // `json:"mtime"`
	InputSamples  []*Sample `json:"input_samples"`
	OutputSamples []*Sample `json:"output_samples"`
	Files         []*File   `json:"files"`
	TemplateID    string    `json:"template_id"`
	TemplateName  string    `json:"template_name"`
}

type Setup struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Attribute  string           `json:"attribute"`
	Properties []*SetupProperty `json:"properties"`
}

type SetupProperty struct {
	ID          string                 `json:"id"`
	Attribute   string                 `json:"attribute"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	OType       string                 `json:"otype"`
	Unit        string                 `json:"unit"`
	Value       map[string]interface{} `json:"value"`
}

type Dataset struct {
}

type Measurement struct {
}

type File struct {
}
