import React from "react";
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";
import AlignItemsList from "./Messages";
import { TextField, Container } from "@mui/material";
import Grid from '@mui/material/Grid';

const Chat = () => {
    const [message, setMessage] = React.useState("")
    const [messages, setMessages] = React.useState([])

    var ws = ChattyService.joinChat()

    ws.onopen = (event) => {
        ws.send(JSON.stringify({
            event: "auth",
            data: UserService.getToken()
        }))
    }
    ws.onmessage = (event) => {
        console.log("from server: " + event.data)
        setMessages((oldMsgs) => {
            return [...oldMsgs, JSON.parse(event.data)]
        })
    }


    const clearMessage = () => {
        setMessage("")
    }

    const sendMessage = () => {
        const msg = {
            event: "message",
            sender: {
                picture: UserService.getParsedToken().picture,
                name: UserService.getParsedToken().name,
            },
            data: message.trim()
        }
        if (msg.data === '') return

        ws.send(JSON.stringify(msg))
        clearMessage()
    }

    return (
        <div>
            <div>
                <Container maxWidth="xl">
                    <Grid container spacing={2}>
                        <Grid item xs={12} style={{
                            height: '95vh',
                            maxHeight: '95vh',
                            display: 'flex',
                            flexDirection: 'column-reverse'
                        }}>
                            <AlignItemsList messages={messages} />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                id="standard-basic"
                                autoComplete="off"
                                label="type message [enter to send]"
                                variant="standard"
                                style={{
                                    width: '100%',
                                }}
                                value={message} onChange={(e) => { setMessage(e.target.value) }}
                                onKeyUp={(e) => {
                                    if (e.key === 'Enter') sendMessage()
                                }}
                            />
                        </Grid>
                    </Grid>
                </Container>
            </div>
        </div>
    )
}

export default Chat;