package main
import(
	"./atm"
	"fmt"
	"./fileHandler"
	)
func main(){
	fmt.Println("Starting server")
	clients := fileHandler.DoStuff();
	i := clients[1].GetCardnr()
	print(i)
	t:= new(atm.Thread)
	t.ReadClients(clients)
	t.ListenForConnections()
}
