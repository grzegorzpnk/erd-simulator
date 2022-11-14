reset

set terminal svg font 'Helvetica,19' size 1200,800

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:100]
set style data histogram
set style histogram rowstack gap 1
set style fill solid border -1

set output 'relocations-rates.svg'

set arrow 1 from -1,20 to 29,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,40 to 29,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,60 to 29,60 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,80 to 29,80 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.8 relative
set xtics rotate by 90 offset 0,-5
set bmargin 7

set ylabel "Relocation Triggering Rate [%]" font 'Helvetica,25'
set xlabel "Time" font 'Helvetica,25'  offset 0,-4
set key inside top left horizontal font "Helvetica, 25" width 1.0

plot newhistogram "M=1" font 'Helvetica,17' offset 0,4.5, \
        'iteration1.dat' using 2:xticlabels(1) title "Cloud Gaming" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=2" font 'Helvetica,17' offset 0,4.5, \
        'iteration2.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "V2x" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=3" font 'Helvetica,17' offset 0,4.5, \
        'iteration3.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "UAV" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=4" font 'Helvetica,17' offset 0,4.5, \
        'iteration4.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars, \
     newhistogram "M=5" font 'Helvetica,17' offset 0,4.5, \
        'iteration5.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
        '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
        '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
        '' using 6:(column(2)+column(3)+column(4)):5 notitle lc rgb "dark-yellow" lw 2 lt 1 with yerrorbars
