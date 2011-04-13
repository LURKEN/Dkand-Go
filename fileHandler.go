package fileHandler 	
import( 
  "os" 
  "fmt" 
  "strings"
  "strconv"
) 
type Client struct{
	cardnr int;
	Saldo chan int
	code int;
	name string;
}

var debug bool  = true;
var clients []Client

func DoStuff() []Client{
	clients := make([]Client,100000);
	clients = ReadWholeFile();
	return clients;
}
//reads clients.txt and inserts values into clients
func ReadWholeFile() []Client{
	var myFile *os.File; 
	var myFileInfo *os.FileInfo;
	var myError os.Error; 

	myFile, myError = os.Open("clients.txt",os.O_RDONLY,0)
	if myError != nil {fmt.Println(myError);}
	myFileInfo,_ = myFile.Stat();	//Stat laser in lite info om filen
	size := myFileInfo.Size;		//Sparar filens storlek 
	WholeFile := make([]byte,size);	//skapar en []byte som ar lika stor som filen.
	
	if debug {fmt.Println("Laser filen!");}
	i,e := myFile.Read(WholeFile);
	if e != nil {
		fmt.Println(i);
		fmt.Println(e);
	}
	WholeFileString := string(WholeFile[0:i]);	//gor om hela filen fran byte[] till en string
	rader := strings.Split(WholeFileString, "\n", -1)	//splitar upp filen pa ny rad.
	l := len(rader);
	clients := make([]Client,l);
	
	if debug {fmt.Println("Hittade dessa Clienter:");}
	
	for i := 0; i < l; i++{	//har kommer det handa andra grejer senare, maste spara varden till clients
		ord := strings.Split(rader[i], " ", 4)	//splittar upp en rad pa mellanrum.

		cardnr, _ := strconv.Atoi(ord[0])
		saldo, _ := strconv.Atoi(ord[1])
		code, _ := strconv.Atoi(ord[2])

		if debug {fmt.Printf("%v %v %v %v %v\n","   ",cardnr,saldo,code, ord[3]);}
		clients[i].cardnr = cardnr;
		s := make(chan int, 1)
		clients[i].Saldo = s
		clients[i].Saldo <- saldo;
		clients[i].code = code;
		clients[i].name = ord[3];
	}
	myFile.Close();
	
	return clients;
}

func WriteToFile(clients []Client){
    f, err := os.Open("clients.txt", os.O_RDWR | os.O_CREATE, 0666) 
    if err != nil {return} 
	if debug {fmt.Println("öppnat filen");}
    defer f.Close() 
	l := len(clients);
	
	for i := 0; i < l; i++ {
		//fmt.Printf("%v %v %v %v\n",clients[i].cardnr,clients[i].	,clients[i].code, clients[i].name);
		cardnr := strconv.Itoa(clients[i].cardnr)
		saldo := strconv.Itoa(<-clients[i].Saldo)
		code := strconv.Itoa(clients[i].code)
		
		_, err = f.Write([]byte(cardnr+" "))
		_, err = f.Write([]byte(saldo+" ")) 
		_, err = f.Write([]byte(code+" ")) 
		_, err = f.Write([]byte(clients[i].name)) 
		if i != l-1 { _, err = f.Write([]byte("\n")) }//Makes it so that the last line doesn't print a \n whitch make it crash the next time you read the file.
	}
	if debug {fmt.Println("Skrivit till filen");}
	f.Close();
	if debug {fmt.Println("Stängt filen");}
}
func (c* Client) GetCardnr() int{ return c.cardnr }
func (c* Client) GetSaldo() chan int{ return c.Saldo }
func (c* Client) GetCode() int{ return c.code }
func (c* Client) GetName() string{ return c.name }
func (c* Client) SetSaldo(in int){ c.Saldo <- in; }
