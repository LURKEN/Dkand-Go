This is a project I(Andreas Starrsjö) and Yuuki Jonsson did at KTH for
our Bachelor degree 2011(http://www.csc.kth.se/utbildning/kth/kurser/DD143X/dkand11/) (Swedish)

atm.go				The Server, it just listens for clients
ATMClient.go		The Client, it connects to the server, you need to configure the IP and PORT in here.
fileHandler.go		The FileHandler, read and writes to clients.txt(where the clients are stored)
main.go				The Main, connects server and filehandler, run this to start the server!

On a windows machine you can rune the Compile(...).bat files to compile the project.
You need to install the Go windows runtime first though.

Any questions? Ask Google. Or maybe maybe me.