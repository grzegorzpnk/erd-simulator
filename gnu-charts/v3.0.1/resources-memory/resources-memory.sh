reset
set terminal postscript eps enhanced colour font 'Helvetica,11'

#set title "MEMORY utilization [per edge type]" font 'Helvetica,25'

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:220]
set style data histogram
set style histogram rowstack gap 2
set style fill solid border -1

set output 'resources-memory.eps'

set arrow 1 from -1,50 to 29,50 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,100 to 29,100 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,150 to 29,150 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
#set arrow 4 from -1,200 to 29,200 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.5 relative
set xtics rotate by 90 offset 0,-5
set bmargin 6

set ylabel "Edge Server MEMORY Utilization [%]" font 'Helvetica,20'

plot newhistogram "M=1", \
       'iteration1.dat' using 2:xticlabels(1) title "local-zone" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "zone" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "international" linecolor rgb "light-gray", \
     newhistogram "M=2", \
       'iteration2.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
     newhistogram "M=3", \
       'iteration3.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
     newhistogram "M=4", \
       'iteration4.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "" linecolor rgb "light-gray", \
     newhistogram "M=5", \
       'iteration5.dat' using 2:xticlabels(1) title "" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "" linecolor rgb "light-gray"
