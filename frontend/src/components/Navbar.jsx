import React from 'react'

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import TerminalIcon from '@mui/icons-material/Terminal';
import Button from '@mui/material/Button';

const Navbar=({handleClick})=>{ 

    return (
      <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
         
         <TerminalIcon sx={{  fontSize: 65 }} style={{marginRight: 20 }}/>
         <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>Práctica 2 - S02</Typography>
         
         <Button color="inherit" style={{fontSize: 20}} onClick={handleClick} >Procesos</Button>
         <Button color="inherit" style={{fontSize: 20}} onClick={handleClick}>Monitor CPU</Button>
         <Button color="inherit" style={{fontSize: 20}} onClick={handleClick}>Monitor RAM</Button>
        </Toolbar>
      </AppBar>
    </Box>
    );
}

export default Navbar ;
