import { useState, useRef } from "react"
import ChattyService from "../services/ChattyService";
import UserService from "../services/UserService";


const Chat = () => {
    const token = UserService.getParsedToken()
    const [message, setMessage] = useState("")

    var ws = useRef(ChattyService.joinChat())

    ws.current.onopen = (event) => {
        ws.current.send(JSON.stringify({
            event: "auth",
            data: UserService.getToken()
        }))
    }
    ws.current.onmessage = (event) => {
        console.log("from server: " + event.data)
    }
    

    const clearMessage = () => {
        setMessage("")
    }

    const sendMessage = () => {
        ws.current.send(JSON.stringify({
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
        </div>
    )
}

export default Chat;