import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';


import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';

import AdministradorTareas from '../components/AdministradorTareas';
const Procesos=({AllTipoProcesos,AllProcesos,AllGenerales})=>{ 

    return (
      <div>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Procesos
            </Typography>
            </center>
         </Paper>

 
         <br/>
         <br/>
         <div style={{marginLeft:100,marginRight:100}}>
            <TableContainer component={Paper} >
               <Table aria-label="simple table" >
                  <TableHead  style={{backgroundColor:"#FFE659"}}>
                     <TableRow >
                        <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Running</b></Typography></TableCell>
                        <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Sleeping </b></Typography></TableCell>
                        <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Stopped </b></Typography></TableCell>
                        <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Zombie </b></Typography></TableCell>
                        <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Total </b></Typography></TableCell>
                     </TableRow>
                  </TableHead>
                  <TableBody>
                     {AllTipoProcesos.map((row,index) => (
                        <TableRow key={`${row.name}${index}`}>
                           <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{row.running}</Typography></TableCell>
                           <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{row.sleeping}</Typography></TableCell>
                           <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{row.stopped}</Typography></TableCell>
                           <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{row.zombie}</Typography></TableCell>
                           <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{row.total}</Typography></TableCell>
                        </TableRow>
                     ))}
                  </TableBody>
               </Table>
            </TableContainer>
         </div>
         
         
         <br/>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Árbol de procesos
            </Typography>
            </center>
         </Paper>
         
         
          <br/>
         <br/>
         <center>
         <AdministradorTareas AllProcesos={AllProcesos} AllGenerales={AllGenerales}/>
         </center>
         <br/>
         <br/>
         <br/>

       

      </div>
    );
}

export default Procesos ;




