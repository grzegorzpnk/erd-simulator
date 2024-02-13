#!/bin/bash

reset

set terminal svg font 'Helvetica,14' size 1400,800

set title "<TODO xLabel>" font 'Helvetica,40' offset 0,-2

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:100]
set style data histogram
set style histogram rowstack gap 2
set style fill solid border -1

set output 'memory-aggregated.svg'

set arrow 1 from -1,20 to 19,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 19,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 19,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 19,80 nohead lc 'red' lw 2  lt 2 dt 1


set boxwidth 0.95 relative
set xtics rotate by 90 offset 0,-6.5
set bmargin 7

set ylabel "Edge Server MEMORY Utilization [%]" font 'Helvetica,20'
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "Optimal-Hybrid" font 'Helvetica,17' offset 0,0.7, \
       'algorithm1.dat' using 2:xticlabels(1) title "City-Level" linecolor rgb "black", \
     newhistogram "EAR-LB" font 'Helvetica,17' offset 0,0.7, \
       'algorithm2.dat' using 2:xticlabels(1) title "Regional-Level" linecolor rgb "dark-gray", \
     newhistogram "EAR-RL" font 'Helvetica,17' offset 0,0.7, \
       'algorithm3.dat' using 2:xticlabels(1) title "Regional-Level" linecolor rgb "dark-gray", \
     newhistogram "EAR-RL-Masked" font 'Helvetica,17' offset 0,0.7, \
       'algorithm4.dat' using 2:xticlabels(1) title "International-Level" linecolor rgb "light-grey"

exit
