import { Avatar, Box, Container, IconButton, InputBase, List, ListItem, ListItemAvatar, ListItemText, Menu, MenuItem, Popover, Toolbar, Tooltip, Typography } from '@mui/material';
import * as React from 'react';
import UserService from '../services/UserService';
import MenuIcon from '@mui/icons-material/Menu';
import MuiAppBar from "@mui/material/AppBar";
import { alpha, styled } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import { debounce, memoize } from 'lodash';

const settings = ['Profile', 'Account', 'Logout'];
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

const Search = styled('div')(({ theme }) => ({
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: alpha(theme.palette.common.white, 0.15),
    '&:hover': {
        backgroundColor: alpha(theme.palette.common.white, 0.25),
    },
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
        marginLeft: theme.spacing(1),
        width: 'auto',
    },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
    padding: theme.spacing(0, 2),
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: 'inherit',
    '& .MuiInputBase-input': {
        padding: theme.spacing(1, 1, 1, 0),
        // vertical padding + font size from searchIcon
        paddingLeft: `calc(1em + ${theme.spacing(4)})`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            width: '12ch',
            '&:focus': {
                width: '20ch',
            },
        },
    },
}));


const Header = (props) => {
    const [anchorElSearch, setAnchorElSearch] = React.useState(null)
    const [anchorElUser, setAnchorElUser] = React.useState(null)
    const [searchedUsers, setSearchedUsers] = React.useState([])
    const entireSearchBar = React.useRef(null)

    const searchRef = React.useRef(null)

    const typingSearchDebounce = React.useMemo(() => memoize(
        debounce((e) => {
            UserService.listUsers(e.target.value).then((res) => {
                setSearchedUsers(res.data)
                setAnchorElSearch(entireSearchBar.current)
                console.log(res.data)
            })
        }, 1000)
    ), [setSearchedUsers, entireSearchBar])

    const handleCloseSearch = () => {
        setAnchorElSearch(null);
    };

    const handleOpenUserMenu = (event) => {
        setAnchorElUser(event.currentTarget);
    };

    const handleCloseUserMenu = () => {
        setAnchorElUser(null);
    };

    const handleChooseOption = (option) => {
        routeOption(option)
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
                        <Box sx={{ flexGrow: 1 }} />
                        <Box sx={{ flexGrow: 1 }}>
                            <Search ref={entireSearchBar}>
                                <SearchIconWrapper>
                                    <SearchIcon />
                                </SearchIconWrapper>
                                <StyledInputBase
                                    ref={searchRef}
                                    placeholder="Search user..."
                                    inputProps={{ 'aria-label': 'search' }}
                                    onChange={typingSearchDebounce}
                                />
                            </Search>
                            <Popover
                                open={Boolean(anchorElSearch)}
                                anchorEl={anchorElSearch}
                                onClose={handleCloseSearch}
                                anchorOrigin={{
                                    vertical: 'bottom',
                                    horizontal: 'left',
                                }}
                                PaperProps={{
                                    style: {
                                        width: `${entireSearchBar.current == null ? 'auto' : entireSearchBar.current.offsetWidth + 'px'}`
                                    },
                                }}
                            >
                                {searchedUsers.length === 0
                                    ?
                                    <Typography sx={{ p: 2 }}>
                                        No matches found
                                    </Typography>
                                    :
                                    <List style={{ overflow: 'auto' }} >
                                        {
                                            searchedUsers.map((user) => {
                                                return (
                                                    <ListItem key={user.username}>
                                                        <ListItemAvatar>
                                                            <Avatar
                                                                imgProps={{ referrerPolicy: "no-referrer" }}
                                                                alt={user.firstName + ' ' + user.lastName}
                                                                src={user.attributes.picture[0]}
                                                            />
                                                        </ListItemAvatar>
                                                        <ListItemText
                                                            sx={{ overflowWrap: 'break-word', wordWrap: 'break-word' }}
                                                            primary={user.firstName + ' ' + user.lastName}
                                                        />
                                                    </ListItem>
                                                )
                                            })
                                        }
                                    </List>


                                }
                            </Popover>
                        </Box>
                        <Box sx={{ flexGrow: 1 }} />
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