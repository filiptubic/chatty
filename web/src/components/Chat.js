import { useState } from "react"
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";


const Chat = () => {
    const token = UserService.getParsedToken()
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
        ws.send(JSON.stringify({
            event: "message",
            data: message
        }))
        console.log('message "' + message + '" sent')
        clearMessage()
    }

    return (
        <div>
            <div>{token.name} ({token.email})</div>
            <div>
                <input value={message} onChange={(e) => { setMessage(e.target.value) }} />
                <button onClick={() => {sendMessage()}}>send</button>
            </div>
            <div>
            {
                messages.map(function(msg, i){
                    if (msg.event !== 'message') return null
                    return <div key={i}>{msg.data}</div>
                })
            }
            </div>
        </div>
    )
}

export default Chat;