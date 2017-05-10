package main

const defaultGrammar = `
{
	"entry": "sentence",
	"rules": {
		"sentence": ["seq", ["ref", "hero"], " ", ["ref", "deed"], " ", ["ref", "villain"]],
		"deed": ["alt", "beat", "bested", "outsmarted"],
		"hero": ["seq", "the ", ["ref", "postrait"], " ", ["alt", "jello", "mayo", "brocolli"]],
		"villain": ["seq", "the ", ["ref", "negtrait"], " ", ["alt", "carrot", "soup", "bean"]],
		"postrait": ["alt", "cunning", "brave", "genuine", "resourceful"],
		"negtrait": ["alt", "treacherous", "reckless", "dumb", "means-justifying"]
	}
}`
