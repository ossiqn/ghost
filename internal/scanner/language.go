package scanner

type Language interface {
	Name() string
	Extensions() []string
	Scan(path string) *Result
	ScanDeps(root string) ([]GhostDep, error)
}

type Registry struct {
	languages map[string]Language
}

func NewRegistry() *Registry {
	return &Registry{
		languages: make(map[string]Language),
	}
}

func (r *Registry) Register(lang Language) {
	for _, ext := range lang.Extensions() {
		r.languages[ext] = lang
	}
}

func (r *Registry) Get(ext string) (Language, bool) {
	lang, ok := r.languages[ext]
	return lang, ok
}

func (r *Registry) All() []Language {
	seen := make(map[string]bool)
	var langs []Language
	for _, lang := range r.languages {
		if !seen[lang.Name()] {
			seen[lang.Name()] = true
			langs = append(langs, lang)
		}
	}
	return langs
}
