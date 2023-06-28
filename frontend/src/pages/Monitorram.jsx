import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { CanvasJSChart } from 'canvasjs-react-charts'

import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';


let dataTemp = [];
const Monitorram = ({ AllGenerales }) => {
   AllGenerales = AllGenerales.length > 0 ? AllGenerales : [{ ramocupada: 0, totalram: 0 }];
   dataTemp.push({ ramocupada: AllGenerales[0].ramocupada, totalram: AllGenerales[0].totalram });

   const data2 = dataTemp.map((item, index) => {
      return { x: index, y: item.ramocupada }
   }).slice(-20);


   const options = {
      animationEnabled: true,
      exportEnabled: true,
      theme: "light2", // "light1", "dark1", "dark2"
      title: {
         text: "Uso de Memoria"
      },
      axisY: {
         title: "MEMORIA",
         suffix: "MB"
      },
      axisX: {
         title: "TIMESTAMP",
         prefix: "T",
         interval: 2
      },
      data: [{
         type: "line",
         toolTipContent: "T {x}: {y} MB",
         dataPoints: data2
      }]
   }

   return (
      <>
         <br />
         <Paper >
            <center>
               <Typography variant="h4" color="inherit" component="div">
                  Monitor RAM
               </Typography>
            </center>
         </Paper>

         <br />
         <br />
         <div>

         </div>


         <div className='centerRam' >
            <CanvasJSChart options={options} />
         </div>


         <br />
         <br />
         <Paper >
            <center>
               <Typography variant="h4" color="inherit" component="div">
                  RAM
               </Typography>
            </center>
         </Paper>

         <br />
         <br />

         <Table aria-label="simple table" >
            <TableHead style={{ backgroundColor: "#FFE659" }}>
               <TableRow >
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>%Usado</b></Typography></TableCell>
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Total </b></Typography></TableCell>
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Consumida </b></Typography></TableCell>
               </TableRow>
            </TableHead>
            <TableBody>
               <TableRow >
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{((AllGenerales[AllGenerales.length - 1].ramocupada) / (AllGenerales[AllGenerales.length - 1].totalram) * 100).toFixed(2)} %</Typography></TableCell>
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{((AllGenerales[AllGenerales.length - 1].totalram)).toFixed(2)} MB</Typography></TableCell>
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div"> {((AllGenerales[AllGenerales.length - 1].ramocupada)).toFixed(2)} MB</Typography></TableCell>
               </TableRow>
            </TableBody>
         </Table>


         <br />
         <br />


      </>
   );
}

export default Monitorram;
