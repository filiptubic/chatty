import React from "react";
import { Box } from '@mui/material';
import './App.css';
import Chat from './components/Chat';
import Header from './components/Header';
import {Sidebar, DrawerHeader} from './components/Sidebar'

function App() {

  const [open, setOpen] = React.useState(false);

  const handleDrawerClose = () => {
    setOpen(false);
  };

  return (
    <div className="App">
      <Box sx={{ display: 'flex' }}>
        <Header
          open={open}
          setOpen={setOpen}
          handleDrawerClose={handleDrawerClose}
          sx={{ flexGrow: 1 }}
        />
        <Sidebar
          sx={{ flexGrow: 1 }}
          open={open}
          setOpen={setOpen}
          handleDrawerClose={handleDrawerClose}
        />
        <Box sx={{ flexGrow: 1 }}>
          <DrawerHeader/>
          <Chat />
        </Box>
      </Box>
    </div>
  );
}

export default App;
