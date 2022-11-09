reset
set terminal postscript eps enhanced colour size 15cm,12cm font 'Helvetica,40'
set boxwidth 0.5
set title "This is chart title" font 'Helvetica,40'

set termoption enhanced
set yrange [0:470]
set style data histograms
set style histogram columnstacked
set style fill pattern  border -1

set output 'sample1.eps'
#set key outside top center
set arrow 1 from -1,225 to 6,225 nohead lc 'midnight-blue' lw 2  lt 2 dt 9
plot 'sample1.dat' using 1, '' using 2, '' using 3, '' using 4, '' using 5, '' using 6:xtic(7)