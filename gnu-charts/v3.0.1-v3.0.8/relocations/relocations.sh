reset
set terminal postscript eps enhanced colour font 'Helvetica,11'

#set title "Number of relocations [per iteration, per app type]" font 'Helvetica,25'

set style line 2 lc rgb 'black' lt 1 lw 1

set termoption enhanced
set yrange [0:50]
set style data histogram
set style histogram rowstack gap 2
set style fill solid border -1

set output 'relocations.eps'

set arrow 1 from -1,10 to 29,10 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 2 from -1,20 to 29,20 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 3 from -1,30 to 29,30 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
set arrow 4 from -1,40 to 29,40 nohead lc 'midnight-blue' lw 2  lt 2 dt 9

set boxwidth 0.5 relative
set xtics rotate by 90 offset 0,-5
set bmargin 6

set ylabel "Number of Successful Relocations [-]" font 'Helvetica,20'

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
