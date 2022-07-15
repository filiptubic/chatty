import * as React from 'react';
import { Fade, List, ListItem, ListItemText, ListItemAvatar, Avatar, Divider } from '@mui/material';

export default function AlignItemsList(prop) {
    const messagesEndRef = React.createRef()

    React.useEffect(() => {
        const scrollToBottom = () => {
            messagesEndRef.current.scrollIntoView({ behavior: 'instant' })
        }
        scrollToBottom()
    }, [prop.messages, messagesEndRef]);

    return (
        <List style={{overflow: 'auto'}} >
            {
                prop.messages.map(function (msg, i) {
                    if (msg.event !== 'message')
                        return null

                    return (
                        <div key={i}>
                            <div>{i > 0 && <Divider variant="inset" component="li" />}</div>
                            <ListItem>
                                <ListItemAvatar>
                                    <Avatar alt={msg.sender.name} imgProps={{ referrerPolicy: "no-referrer" }} src={msg.sender.picture} />
                                </ListItemAvatar>
                                <Fade in={true}>
                                    <ListItemText style={{ overflowWrap: 'break-word', wordWrap: 'break-word' }} primary={msg.data} secondary={msg.sender.name} />
                                </Fade>
                            </ListItem>
                        </div>
                    )
                })
            }
            <ListItem>
                <div ref={messagesEndRef} />
            </ListItem>
        </List>
    );
}
