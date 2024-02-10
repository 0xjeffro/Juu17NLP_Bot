package utils

import (
	"Juu17NLP_Bot/orm"
	"regexp"
	"strings"
)

func KeywordsAnalysis(text string) []string {
	db := orm.GetConn()
	var rules []orm.Rules
	db.Limit(300).Where(&orm.Rules{Type: "keyword"}).Find(&rules)

	keywords := make([]string, 0)
	for _, rule := range rules {
		keyword := rule.Content
		if strings.Contains(text, keyword) {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func RegularExpressionAnalysis(text string) []string {
	db := orm.GetConn()
	var rules []orm.Rules
	db.Limit(300).Where(&orm.Rules{Type: "regex"}).Find(&rules)

	regex := make([]string, 0)
	for _, rule := range rules {
		keyword := rule.Content
		found, _ := regexp.MatchString(keyword, text)
		if found {
			regex = append(regex, keyword)
		}
	}
	return regex
}
