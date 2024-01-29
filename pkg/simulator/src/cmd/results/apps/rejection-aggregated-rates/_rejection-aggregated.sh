#!/bin/bash

reset

set terminal svg font 'Helvetica,19' size 1400,800

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:20]
set style data histogram
set style histogram rowstack gap 1
set style fill solid border -1

set output 'rejection-rates.svg'

set arrow 1 from -1,20 to 11,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 11,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 11,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 11,80 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.95 relative
set xtics rotate by 90 offset 0,-5
set bmargin 7

set ylabel "REJECTION Rate [%]" font 'Helvetica,25'
set xlabel "Time" font 'Helvetica,25'  offset 0,-4
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "" font 'Helvetica,17' offset 0,4.5, \
        'rejection.dat' using 2:xticlabels(1) title "Cloud Gaming" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "V2x" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "UAV" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars
