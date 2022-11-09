set terminal jpeg giant font "Helvetica" 16

set output 'stacked-histogram-3.jpg'
set key left

set style line 2 lc rgb 'black' lt 1 lw 1
set grid y
set style data histograms
set style histogram rowstacked
set boxwidth 0.5
set style fill pattern border -1
set ytics 10 nomirror
set yrange [:40]
set ylabel "Number of failed search"
set ytics 10
set key inside top center horizontal


plot 'stacked-histogram-3.dat' using 2:xtic(1) t "1" ls 2, '' using 3 t "2" ls 2, '' using 4 t "3" ls 2, '' using 5 t "4" ls 2, '' using 6 t "5" ls 2
