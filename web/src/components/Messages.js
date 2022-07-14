import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Divider from '@mui/material/Divider';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';

export default function AlignItemsList(prop) {

    return (
        <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
            {
                prop.messages.map(function (msg, i) {
                    if (msg.event !== 'message')
                        return null

                    return (
                        <div key={i}>
                            <div>{i > 0 && <Divider variant="inset" component="li" />}</div>
                            <ListItem alignItems="flex-start">
                                <ListItemAvatar>
                                    <Avatar alt={msg.sender.name} imgProps={{referrerPolicy: "no-referrer"}} src={msg.sender.picture} />
                                </ListItemAvatar>
                                <ListItemText
                                    primary={msg.data}
                                />
                            </ListItem>
                        </div>
                    )
                })
            }
        </List>
    );
}
