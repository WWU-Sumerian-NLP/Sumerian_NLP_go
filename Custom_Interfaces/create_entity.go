package custom_interfaces

type CustomEntity struct {
	Name       string
	Tag        string
	PathToList string
}

//attributes are extracted afterwards not at labeling

func newCustomEntity(name string, tag string, pathToList string) *CustomEntity {
	return &CustomEntity{Name: name, Tag: tag, PathToList: pathToList}
}

/*

Give entity_extractor pipeline class a list of Custom Entity objects

Iterate through each of them and append their pathToList to the NER list
The file should already label as so

name, entity_tag, translations?

In the future, the entity tag may not be labeled

The name may seem redundant, but its for debugging purposes and to be later exposed to the user


*/
