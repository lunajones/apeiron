package politic

import (
	"github.com/lunajones/apeiron/service/politic/consts"
)

var CityRelations = map[string]map[string]consts.RelationStatus{
	"A": {
		"B": consts.Hostile,
		"C": consts.Friendly,
	},
	"B": {
		"A": consts.Hostile,
		"C": consts.Neutral,
	},
}

func FactionsAreHostile(cityA, cityB string) bool {
	return CityRelations[cityA][cityB] == consts.Hostile
}

func SetRelationStatus(cityA, cityB string, status consts.RelationStatus) {
	if _, ok := CityRelations[cityA]; !ok {
		CityRelations[cityA] = make(map[string]consts.RelationStatus)
	}
	CityRelations[cityA][cityB] = status
}
