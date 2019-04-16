This is an implementation of a Markov chain text generator in go-lang. It uses
a map to store the Markov chains found in the input file, and generates an
output file using their relative probability of concurrence to the previous
n-gram.

It is my first time writing a fully fledged program in go and it works, but I
feel like I really just wrote python with different syntax. I was trying to
avoid doing that, but I couldn't help it - it's what I am more familiar with.

If go is installed on the machine, it can be run with something like:
    go run ./markov.go {input filename} {order of n-gram} {output filename} \
           {length generated}

    for example:
    go run ./markov.go alice.txt 5 test.txt 200
