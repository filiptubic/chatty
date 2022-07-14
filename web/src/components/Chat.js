import { useState } from "react"
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";
import AlignItemsList from "./Messages";
import { TextField, Button, Container } from "@mui/material";

const Chat = () => {
    const [message, setMessage] = useState("")
    const [messages, setMessages] = useState([])

    var ws = ChattyService.joinChat()

    ws.onopen = (event) => {
        ws.send(JSON.stringify({
            event: "auth",
            data: UserService.getToken()
        }))
    }
    ws.onmessage = (event) => {
        console.log("from server: " + event.data)
        setMessages([...messages, JSON.parse(event.data)])
    }


    const clearMessage = () => {
        setMessage("")
    }

    const sendMessage = () => {
        console.log(UserService.getParsedToken())
        const msg = {
            event: "message",
            sender: {
                picture: UserService.getParsedToken().picture,
                name: UserService.getParsedToken().name,
            },
            data: message
        }
        ws.send(JSON.stringify(msg))
        console.log('message "' + message + '" sent')
        clearMessage()
    }

    return (
        <div>
            {/* <div>{token.name} ({token.email})</div>
            <div>
                <input value={message} onChange={(e) => { setMessage(e.target.value) }} />
                <button onClick={() => { sendMessage() }}>send</button>
            </div> */}
            <div>
                <Container maxWidth="sm">
                    <div>
                        <TextField id="standard-basic" label="enter message" variant="standard" value={message} onChange={(e) => { setMessage(e.target.value) }} />
                        <Button onClick={() => { sendMessage() }} variant="contained">
                            Send
                        </Button>
                    </div>
                    <AlignItemsList messages={messages} />
                </Container>
            </div>
        </div>
    )
}

export default Chat;