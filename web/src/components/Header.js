import { Avatar, Box, Container, IconButton, Menu, MenuItem, Toolbar, Tooltip, Typography } from '@mui/material';
import * as React from 'react';
import UserService from '../services/UserService';
import MenuIcon from '@mui/icons-material/Menu';
import MuiAppBar from "@mui/material/AppBar";
import { styled } from '@mui/material/styles';

const settings = ['Profile', 'Account', 'Logout'];
const drawerWidth = 240;

const Header = (props) => {
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

    var open = props.open
    var setOpen = props.setOpen

    const handleDrawerOpen = () => {
        setOpen(true);
    };

    const AppBar = styled(MuiAppBar, {
        shouldForwardProp: (prop) => prop !== "open"
      })(({ theme, open }) => ({
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(["width", "margin"], {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.leavingScreen
        }),
        ...(open && {
          marginLeft: drawerWidth,
          width: `calc(100% - ${drawerWidth}px)`,
          transition: theme.transitions.create(["width", "margin"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen
          })
        })
      }));


    return (
        <Box sx={{display: 'flex'}}>
            <AppBar position="fixed" open={open}>
                <Container maxWidth={false}>
                    <Toolbar disableGutters>
                        <Box sx={{ flexGrow: 0 }}>
                            <IconButton
                                color="inherit"
                                aria-label="open drawer"
                                onClick={handleDrawerOpen}
                                edge="start"
                                sx={{
                                    marginRight: 5,
                                    ...(open && { display: 'none' }),
                                }}
                            >
                                <MenuIcon />
                            </IconButton>
                        </Box>
                        <Box sx={{ flexGrow: 1 }}></Box>
                        <Box sx={{ flexGrow: 0 }}>
                            <Tooltip title="Open settings">
                                <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                                    <Avatar imgProps={{ referrerPolicy: "no-referrer" }} alt={UserService.getParsedToken().name} src={UserService.getParsedToken().picture} />
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
                                    <MenuItem key={setting} onClick={() => { handleChooseOption(setting) }}>
                                        <Typography textAlign="center">{setting}</Typography>
                                    </MenuItem>
                                ))}
                            </Menu>
                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>
        </Box>
    )
}

export default Header