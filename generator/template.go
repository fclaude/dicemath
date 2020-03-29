package generator

import "math/rand"

type exercise struct {
	A int
	B int
}

func genExercise() exercise {
	return exercise{A: rand.Intn(6) + 1, B: rand.Intn(6) + 1}
}

const N = 5

type exerciseGroup struct {
	Group1 []exercise
	Group2 []exercise
	Group3 []exercise
	Group4 []exercise
	Name   string
}

func generateGroups(name string) exerciseGroup {
	seen := make(map[exercise]bool, 36)
	randomGroup := func() []exercise {
		result := make([]exercise, N)
		for i := range result {
			ex := genExercise()
			for seen[ex] {
				ex = genExercise()
			}
			seen[ex] = true
			result[i] = ex
		}
		return result
	}
	return exerciseGroup{
		Name:   name,
		Group1: randomGroup(),
		Group2: randomGroup(),
		Group3: randomGroup(),
		Group4: randomGroup(),
	}
}

const sheetTemplate = `
\documentclass[12pt,letter]{article}

\usepackage{fullpage}
\usepackage{amsmath}
\usepackage{epsdice}
\usepackage{mathptmx}
\usepackage{anyfontsize}

\title{ [[.Name]] worksheet}
\author{}
\date{}
\begin{document}


\newcommand{\dice}[1]{{\fontsize{50}{50}\selectfont \epsdice{#1}}}
\newcommand{\resultline}{..........}
\newcommand{\childsum}[2]{\begin{equation*}{\LARGE \dice{#1}~~\raisebox{14pt}{+}~~\dice{#2}~~\raisebox{14pt}{=}~~\resultline}\end{equation*}}

\twocolumn
\maketitle
\thispagestyle{empty}

[[range .Group1]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group2]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group3]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group4]]\childsum{[[.A]]}{[[.B]]}
[[end]]

\end{document}
`
