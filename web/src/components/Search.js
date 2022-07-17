import React from "react"
import UserService from "../services/UserService"
import { Avatar,InputBase, List, ListItem, ListItemAvatar, ListItemText, Popover, Typography } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { debounce, memoize } from 'lodash';
import { alpha, styled } from '@mui/material/styles';
import ChattyService from '../services/ChattyService'
import { useNavigate } from "react-router-dom";

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


const SearchUsers = () => {
    const navigate = useNavigate();
    const entireSearchBar = React.useRef(null)
    const searchRef = React.useRef(null)
    const [anchorElSearch, setAnchorElSearch] = React.useState(null)
    const [searchedUsers, setSearchedUsers] = React.useState([])

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

    const handleOnClickUser = (user) => {
        ChattyService.createChat(user.id).then((res)=>{
            console.log(res)
            navigate(`/${res.data}`);
        })
        handleCloseSearch()
    }

    return (
        <div>
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
                                    <ListItem key={user.username} onClick={() => handleOnClickUser(user)}>
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
        </div>
    )
}

export default SearchUsers