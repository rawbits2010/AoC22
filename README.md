# AoC22
My Advent of Code 2022 competition solutions written in Go.

For all kinds of legal reasons, I won't provide the actual puzzles or input data. You can get those by visiting the competition's webite [adventofcode.com](https://adventofcode.com/2022).

## Usage

Just build the solution corresponding to the day you want from the 'cmd' directory.

You can provide input using any of the 3 options - implemented by the inputhandler package in the 'internal' directory.

First method is to use the commandline, with the input lines separated by a ';'. For example you can run 'day02' with:

`./day02 -p "A X;A Y;A Z"`

Second method is to create a text file containing the actual input lines and provide it using it's path like so:

`./day02 -f input.txt`

The third method is the provide an URL that will do a GET request to fetch the input data:

`./day02 -w https://adventofcode.com/2022/day/2/input`

NOTE: If you want to use this method with the Advent of Code site like above, you need to login to the site and provide your 'session' cookie value in the 'session.txt' file.

## Legal stuff

These solutions are provided as is and I don't take any responsibility for what they cause. Use at your own risk!

You are free to compile and study the code in this repository. Representing it as your own is strictly forbidden and generally rude. This applies to the whole or any parts of it.

(Not that anything in here is particularly valuable or usable outside of the competition. Also, with the different solutions, I aim to learn different aspects of the language.)

Enjoy!