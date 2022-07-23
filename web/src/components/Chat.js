import React from "react";
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";
import AlignItemsList from "./Messages";
import { TextField, Container, Typography, Zoom } from "@mui/material";
import Grid from '@mui/material/Grid';
import { Box } from "@mui/system";
import debounce from "lodash/debounce";
import { useParams } from "react-router-dom";

const Chat = () => {
    const params = useParams();
    const [message, setMessage] = React.useState('')
    const [messages, setMessages] = React.useState([])
    const [typing, setTyping] = React.useState('')
    const [typingShowed, setTypingShowed] = React.useState(false)
    const typingRecvDebounce = React.useRef(debounce(() => {
        setTyping('')
        setTypingShowed(false)
    }, 500))
    const ws = React.useRef(null)

    React.useEffect(()=>{
        setMessages([])
        setMessage('')
        setTyping('')
        setTypingShowed(false)
        ws.current = ChattyService.joinChat(params.chatId)
    },[params.chatId])

    const typingSendDebounce = React.useRef(debounce(() => {
        ws.current.send(JSON.stringify({
            event: "typing",
            sender: {
                picture: UserService.getParsedToken().picture,
                name: UserService.getParsedToken().name,
                username: UserService.getParsedToken().preferred_username
            }
        }))
    }, 100))

    if (ws.current != null) {
        ws.current.onopen = (event) => {
            ws.current.send(JSON.stringify({
                event: "auth",
                data: UserService.getToken()
            }))
        }
    }

    if (ws.current != null) {
        ws.current.onmessage = (event) => {
            const msg = JSON.parse(event.data)
            console.log(msg)
            if (msg.event === 'message') {
                setMessages((oldMsgs) => {
                    return [...oldMsgs, msg]
                })
            } else if (msg.event === 'typing') {
                typingRecvDebounce.current()
                if (typingShowed) return
                if (msg.sender.username === UserService.getParsedToken().preferred_username) return

                setTyping(msg.sender.name + ' typing...')
                setTypingShowed(true)
            }
        }
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

        ws.current.send(JSON.stringify(msg))
        clearMessage()
    }

    const onTyping = (e) => {
        typingSendDebounce.current()
        setMessage(e.target.value)
    }

    return (
        <Container maxWidth="xl" sx={{flexGrow: 1}}>
            <Grid container spacing={2}>
                <Grid item xs={12} style={{
                    height: '85vh',
                    maxHeight: '85vh',
                    display: 'flex',
                    flexDirection: 'column-reverse'
                }}>
                    <AlignItemsList messages={messages} />
                </Grid>
                <Grid item xs={2}>
                    <Box sx={{ display: 'flex', height: "20px" }}>
                        <Zoom in={typingShowed}>
                            <Typography variant="body2"><i>{typing}</i></Typography>
                        </Zoom>
                    </Box>
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
                        value={message}
                        onChange={onTyping}
                        onKeyUp={(e) => {
                            if (e.key === 'Enter') sendMessage()
                        }}
                    />
                </Grid>
            </Grid>
        </Container>
    )
}

export default Chat;