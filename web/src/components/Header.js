import { Box, Container, IconButton, Toolbar } from '@mui/material';
import * as React from 'react';
import MenuIcon from '@mui/icons-material/Menu';
import MuiAppBar from "@mui/material/AppBar";
import { styled } from '@mui/material/styles';
import Search from './Search';
import Settings from './Settings';

const drawerWidth = 240;

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

const Header = (props) => {
    var open = props.open
    var setOpen = props.setOpen

    const handleDrawerOpen = () => {
        setOpen(true);
    };

    return (
        <Box sx={{ display: 'flex' }}>
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
                        <Box sx={{ flexGrow: 1 }}/>
                        <Box sx={{ flexGrow: 1 }}>
                            <Search/>
                        </Box>
                        <Box sx={{ flexGrow: 1 }}/>
                        <Box sx={{ flexGrow: 0 }}>
                            <Settings/>
                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>
        </Box>
    )
}

export default Header