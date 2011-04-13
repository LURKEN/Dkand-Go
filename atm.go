package atm
import(
	"net"
	"fmt"
	"os"
	"bufio"
	"strconv"
	"./fileHandler"
)
var clients []fileHandler.Client
type Thread struct{
	listener *net.TCPListener
	login chan int
}
type Cl struct{
	con net.Conn
	buf []byte
	client *fileHandler.Client
}
func (t *Thread) run(){}
func (t *Thread) ReadClients(cs []fileHandler.Client){
	clients = cs
}
func (t *Thread) ListenForConnections(){
	ip := net.ParseIP("127.0.0.1")
	addr := &net.TCPAddr{ip, 9999};
	listener, err := net.ListenTCP("tcp", addr);
	if err != nil {fmt.Println(err)}
	go t.Command()
	login := make(chan int, 1)
	t.login = login
	t.listener = listener
	t.login<-1
	for{
		<-t.login
		go t.AcceptClient()
	}
}
func (t *Thread) AcceptClient(){
	fmt.Println("Listening for clients")
	conn, _ := t.listener.AcceptTCP()
	t.login<-1
	c := new(Cl)
	c.con = conn
	buf := make([]byte,1024)
	c.buf = buf
	length := len(clients)
	for i := 0; i < 3 ; i++ {	//tre försök att logga in
		c.con.Write([]byte("Cardnr:"))
		size, _ := c.con.Read(c.buf)//waiting for cardnr
		size = size-1; // ta bort \n i slutet
		size = 4;
		if size > 0 {
			for j:=0;j<length;j++{
				tmp := clients[j].GetCardnr()//tar ett kortnummer från de sparade
				tmp2, _ := strconv.Atoi(string(c.buf[0:size]))//skriv in ett kortnummer		
				if tmp == tmp2 {
					fmt.Println("hittat kund, väntar på passwd")
					c.con.Write([]byte("Code:"))	
					nr2, err := c.con.Read(c.buf)//waiting for code
					nr2 = 1;
					if err != nil {print(" error")}
					if(nr2 > 0){
						code := clients[j].GetCode()
						code2, _ := strconv.Atoi(string(c.buf[0:(nr2)]))
						if code == code2 {c.loggedIn(j)}
					}
				}
			}
			fmt.Println("WRONG INPUT!")
		}
		c.con.Write([]byte("error"))//didnt find cardnr
	}
}

func (c *Cl) loggedIn(kundID int){
	fmt.Println("KUNDEN ÄR NU INLOGGAD!")
	c.con.Write([]byte("loggedin"))	//viktig! Klienten checkar efter exakt loggedin för att logga in
	tmp,_ := c.con.Read(c.buf) //handskakning
	if tmp > 0 {}
	for {
		c.PrintMeny(kundID);
		size, _:= c.con.Read(c.buf)//waiting for input
		input, _ := strconv.Atoi(string(c.buf[0:(size-2)]))
		if input == 1 {
				c.con.Write([]byte("remove how much?"))//write
				nr, _ := c.con.Read(c.buf)//read
				antal, _ := strconv.Atoi(string(c.buf[0:nr-2]))
				c.removeMoney(kundID,antal)	//write
		}
		if input == 2 {
				c.con.Write([]byte("insert how much?"))
				nr2, _ := c.con.Read(c.buf)
				antal, _ := strconv.Atoi(string(c.buf[0:(nr2-2)]))	
				c.insertMoney(kundID,antal)	
		}
		if input == 3  {
				fmt.Println("logout")
				c.con.Write([]byte("logout\n"))
				logout()
				break;
		}
		if input == 4 {
			fmt.Println("EXIT")
			c.con.Write([]byte("EXIT\n"))
			exit()
		}
	}
}
func (c *Cl) PrintMeny(kundID int){
	saldo:=<-clients[kundID].Saldo
	CardNr := strconv.Itoa(clients[kundID].GetCardnr())
	str:=strconv.Itoa(saldo)
	c.con.Write([]byte("CardNr: "+CardNr+" Saldo: "+str +" \n1. removeMoney\n2. insertMoney\n3. logout\n4. exit\n"))
	clients[kundID].Saldo<-saldo
}
func (c *Cl) printBalance(kundID int){
	saldo:=<-clients[kundID].Saldo
	str:=strconv.Itoa(saldo)
	c.con.Write([]byte(str))
	clients[kundID].Saldo<-saldo
}
func (c *Cl) removeMoney(kundId int, antal int){
	saldo := <- clients[kundId].Saldo
	clients[kundId].Saldo <- saldo - antal
}
func (c *Cl) insertMoney(kundId int, antal int){
	saldo :=<-clients[kundId].Saldo
	clients[kundId].Saldo<- saldo + antal
}

func exit(){os.Exit(0)}
func logout(){}
func (t *Thread) Command(){
	reader := bufio.NewReader(os.Stdin)
	command, _ := reader.ReadString('\n')
	fmt.Println(command)
	e := t.listener.Close()
	fileHandler.WriteToFile(clients);
	print(e)
	os.Exit(0)
}
