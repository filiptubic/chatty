import { useState } from "react"
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";
import AlignItemsList from "./Messages";
import { TextField, Button, Container } from "@mui/material";
import Grid from '@mui/material/Grid';
import SendIcon from '@mui/icons-material/Send';

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
            <div>
                <Container maxWidth="sm">
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <TextField
                                id="standard-basic"
                                label="enter message"
                                variant="standard"
                                style={{ width: '100%' }}
                                value={message} onChange={(e) => { setMessage(e.target.value) }} />
                        </Grid>
                        <Grid item xs={12}>
                            <Button onClick={() => { sendMessage() }} style={{ width: '100%' }} variant="contained">
                                <span style={{ paddingRight: "10px" }}>Send</span> <SendIcon />
                            </Button>
                        </Grid>
                        <Grid item xs={12}>
                            <AlignItemsList messages={messages} />
                        </Grid>
                    </Grid>

                </Container>
            </div>
        </div>
    )
}

export default Chat;