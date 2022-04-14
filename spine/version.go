package spine

type SpecificationVersionDataType SpecificationVersionType

type SpecificationVersionDataElementsType struct{}

type SpecificationVersionListDataType struct {
	SpecificationVersionData []SpecificationVersionDataType `json:"specificationVersionData,omitempty"`
}

type SpecificationVersionListDataSelectorsType struct{}
