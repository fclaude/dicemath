package generator

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
\newcommand{\childsum}[2]{\begin{equation*}{\LARGE \dice{#1}~~\raisebox{14pt}{[[.Operation]]}~~\dice{#2}~~\raisebox{14pt}{=}~~\resultline}\end{equation*}}

\twocolumn
\maketitle
\thispagestyle{empty}

[[range .Group1]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group2]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[.Separator]]

[[range .Group3]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group4]]\childsum{[[.A]]}{[[.B]]}
[[end]]

\end{document}
`
