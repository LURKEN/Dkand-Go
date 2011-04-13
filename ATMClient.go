package main
import(
	"net"
	"fmt"
	"bufio"
	"os"
)
type Client struct{
	con net.Conn
	buff []byte
	reader *bufio.Reader
}
func main(){
	c := new(Client)
	c.run()
}

func (c *Client) run(){
	buff := make([]byte, 1000)
	c.buff = buff
	c.con = c.connect()
	reader := bufio.NewReader(os.Stdin)
	c.reader = reader
	c.login()
}

func (c *Client) connect() (conn net.Conn){
	conn, err := net.Dial("tcp4", "", "localhost:9999")
	if err != nil {fmt.Println(err)}
	return conn
}

func (c *Client) login() {
	for i:=0;i<3;i++{
		c.recieve()//ta emot card-text
		c.sendCommand()
		str := c.recieve2()//ta emot passwd-text
		if str == "error"{fmt.Println("FEL KORTNR!")
		}else{
			fmt.Println(str)
			c.sendCommand()
			//check if accepted
			msg := c.recieve2()
			c.con.Write([]byte("handskakning"))
			if msg == "loggedin" {//important!
				c.loggedIn()
			}else{fmt.Println("FEL KOD!")}
		}
	}
	fmt.Println("exit")
}

func (c *Client) loggedIn(){
	for{
		c.recieve()
		c.sendCommand()
	}
}

func (c *Client) sendCommand(){
	command, e := c.reader.ReadString('\n')
	if e != nil {print(e)}
	c.con.Write([]byte(command))
}

func (c *Client) recieve(){
	nr, err := c.con.Read(c.buff)
	if err != nil {print(err)}
	if nr > 0 {fmt.Println(string(c.buff[0:nr]))}
}

func (c *Client) recieve2() (string){
	nr, err := c.con.Read(c.buff)
	if err != nil {print(err)}
	if nr > 0 {return string(c.buff[0:nr])}
	return "error";
}