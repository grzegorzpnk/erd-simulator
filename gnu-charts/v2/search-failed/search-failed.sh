reset
set terminal postscript eps enhanced colour font 'Helvetica,11'

#set title "Number of failed search [per iteration, per app type]" font 'Helvetica,25'

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:7]
set style data histogram
set style histogram rowstack gap 2
set style fill solid border -1

set output 'search-failed.eps'

set arrow 1 from -1,1 to 29,1 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,2 to 29,2 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,3 to 29,3 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,4 to 29,4 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 5 from -1,5 to 29,5 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.5 relative
set xtics rotate by 90 offset 0,-5
set bmargin 6

set ylabel "Number of Failed Search [-]" font 'Helvetica,20'

plot newhistogram "M=1", \
       'iteration1.dat' using 2:xticlabels(1) title "cg" linecolor rgb "black", \
                                 '' using 3:xticlabels(1) title "v2x" linecolor rgb "dark-gray", \
                                 '' using 4:xticlabels(1) title "uav" linecolor rgb "light-gray", \
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
