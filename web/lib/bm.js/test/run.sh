while true; do
	echo "open http://localhost:8000/test/test.html"
	find test -type f -print | entr -rd python -m http.server
	echo "hit ^C again to quit" && sleep 1
done

