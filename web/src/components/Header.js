import { AppBar, Avatar, Box, Container, IconButton, Menu, MenuItem, Toolbar, Tooltip, Typography } from '@mui/material';
import * as React from 'react';
import UserService from '../services/UserService';

const settings = ['Profile', 'Account', 'Logout'];

const Header = () => {
    const [anchorElUser, setAnchorElUser] = React.useState(null);

    const handleOpenUserMenu = (event) => {
        setAnchorElUser(event.currentTarget);
    };

    const handleCloseUserMenu = () => {
        setAnchorElUser(null);
    };

    const handleChooseOption = (option) => {
        routeOption(option)
        console.log(option)
        handleCloseUserMenu()
    };

    const routeOption = (option) => {
        if (option === 'Logout') {
            UserService.doLogout()
        }
    }

    return (
        <AppBar position="static">
            <Container maxWidth="xl">
                <Toolbar disableGutters>
                    <Box sx={{flexGrow: 1}}></Box>
                    <Box sx={{ flexGrow: 0 }}>
                        <Tooltip title="Open settings">
                            <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                                <Avatar alt={UserService.getParsedToken().name} src={UserService.getParsedToken().picture} />
                            </IconButton>
                        </Tooltip>
                        <Menu
                            sx={{ mt: '45px' }}
                            id="menu-appbar"
                            anchorEl={anchorElUser}
                            anchorOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            open={Boolean(anchorElUser)}
                            onClose={handleCloseUserMenu}
                        >
                            {settings.map((setting) => (
                                <MenuItem key={setting} onClick={() => {handleChooseOption(setting)}}>
                                    <Typography textAlign="center">{setting}</Typography>
                                </MenuItem>
                            ))}
                        </Menu>
                    </Box>
                </Toolbar>
            </Container>
        </AppBar>
    )
}

export default Header