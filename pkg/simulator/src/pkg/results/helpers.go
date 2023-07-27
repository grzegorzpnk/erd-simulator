package results

import (
	"fmt"
)

type percentageCounter struct {
	total   float64
	divisor float64
}

func (pc *percentageCounter) getPercentage() float64 {
	return pc.total / pc.divisor * 100
}

func initializeEmpty2DArray() [][]float64 {
	a := make([][]float64, 5)
	for i := range a {
		a[i] = make([]float64, 4)
	}
	return a
}

func iteratorConstraint(iter int) bool {
	if iter >= 5 {
		return true
	}
	return false
}

func generateRatesScript(ratesType, xLabel, yLabel string, iterFiles []string) string {

	script := fmt.Sprintf(`#!/bin/bash

reset

set terminal svg font 'Helvetica,19' size 1200,800

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:100]
set style data histogram
set style histogram rowstack gap 1
set style fill solid border -1

set output '%s-rates.svg'

set arrow 1 from -1,20 to 29,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 29,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 29,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 29,80 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.8 relative
set xtics rotate by 90 offset 0,-5
set bmargin 7

set ylabel "%s" font 'Helvetica,25'
set xlabel "%s" font 'Helvetica,25'  offset 0,-4
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "M=1" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "Cloud Gaming" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=2" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "V2x" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=3" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "UAV" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=4" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=5" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars
`, ratesType, yLabel, xLabel, iterFiles[0], iterFiles[1], iterFiles[2], iterFiles[3], iterFiles[4])

	return script
}

func generateAggregatedRatesScriptApps(ratesType, xLabel, yLabel string, iterFile string) (script string) {
	fmt.Println("test")
	if ratesType == "rejection" {
		script = fmt.Sprintf(`#!/bin/bash

reset

set terminal svg font 'Helvetica,19' size 1200,800

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:20]
set style data histogram
set style histogram rowstack gap 1
set style fill solid border -1

set output '%s-rates.svg'

set arrow 1 from -1,20 to 11,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 11,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 11,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 11,80 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.95 relative
set xtics rotate by 90 offset 0,-5
set bmargin 7

set ylabel "%s" font 'Helvetica,25'
set xlabel "%s" font 'Helvetica,25'  offset 0,-4
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "Cloud Gaming" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "V2x" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "UAV" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars
`, ratesType, yLabel, xLabel, iterFile)
	} else {
		script = fmt.Sprintf(`#!/bin/bash

reset

set terminal svg font 'Helvetica,19' size 1200,800

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:100]
set style data histogram
set style histogram rowstack gap 1
set style fill solid border -1

set output '%s-rates.svg'

set arrow 1 from -1,20 to 11,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 11,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 11,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 11,80 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.95 relative
set xtics rotate by 90 offset 0,-5
set bmargin 7

set ylabel "%s" font 'Helvetica,25'
set xlabel "%s" font 'Helvetica,25'  offset 0,-4
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "" font 'Helvetica,17' offset 0,4.5, \
        '%s' using 2:xticlabels(1) title "Cloud Gaming" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "V2x" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "UAV" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars
`, ratesType, yLabel, xLabel, iterFile)
	}

	return
}

func generateAggregatedRatesScriptMecs(resType, title, yLabel string, iterFiles []string) string {

	script := fmt.Sprintf(`#!/bin/bash

reset

set terminal svg font 'Helvetica,14' size 1200,800

set title "%s" font 'Helvetica,40' offset 0,-2

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:100]
set style data histogram
set style histogram rowstack gap 2
set style fill solid border -1

set output '%s-aggregated.svg'

set arrow 1 from -1,20 to 19,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 19,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 19,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 19,80 nohead lc 'red' lw 2  lt 2 dt 1


set boxwidth 0.95 relative
set xtics rotate by 90 offset 0,-6.5
set bmargin 7

set ylabel "%s" font 'Helvetica,20'
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "Optimal-Hybrid" font 'Helvetica,17' offset 0,0.7, \
       '%s' using 2:xticlabels(1) title "City-Level" linecolor rgb "black", \
     newhistogram "Heuristic-Hybrid" font 'Helvetica,17' offset 0,0.7, \
       '%s' using 2:xticlabels(1) title "Regional-Level" linecolor rgb "dark-gray", \
     newhistogram "Heuristic-EAR" font 'Helvetica,17' offset 0,0.7, \
       '%s' using 2:xticlabels(1) title "Regional-Level" linecolor rgb "dark-gray", \
     newhistogram "ML-Masked" font 'Helvetica,17' offset 0,0.7, \
       '%s' using 2:xticlabels(1) title "International-Level" linecolor rgb "light-grey", \
     newhistogram "ML-NonMasked" font 'Helvetica,17' offset 0,0.7, \
       '%s' using 2:xticlabels(1) title "" linecolor rgb "black"

exit
`, title, resType, yLabel, iterFiles[0], iterFiles[1], iterFiles[2], iterFiles[3], iterFiles[4])

	return script
}
