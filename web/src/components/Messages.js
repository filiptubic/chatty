import * as React from 'react';
import { List, ListItem, ListItemText, ListItemAvatar, Avatar, Divider } from '@mui/material';

export default function AlignItemsList(prop) {

    return (
        <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
            {
                prop.messages.map(function (msg, i) {
                    if (msg.event !== 'message')
                        return null

                    return (
                        <div key={i}>
                            <div>{i > 0 && <Divider variant="inset" component="li" />}</div>
                            <ListItem alignItems="flex-start">
                                <ListItemAvatar>
                                    <Avatar alt={msg.sender.name} imgProps={{ referrerPolicy: "no-referrer" }} src={msg.sender.picture} />
                                </ListItemAvatar>
                                <ListItemText style={{ overflowWrap: 'break-word', wordWrap: 'break-word' }} primary={msg.data} />
                            </ListItem>
                        </div>
                    )
                })
            }
        </List>
    );
}
// overflow-wrap: break-word;
// word-wrap: break-word;