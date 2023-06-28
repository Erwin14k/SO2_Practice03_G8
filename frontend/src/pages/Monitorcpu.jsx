import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { CanvasJSChart } from 'canvasjs-react-charts'

   let dataTemp = [];

const Monitorcpu=({AllGenerales})=>{ 
   AllGenerales = AllGenerales.length > 0 ? AllGenerales : [{totalcpu:0}];
   dataTemp.push({totalcpu:AllGenerales[0].totalcpu});

   const data2 = dataTemp.map((item, index) => {
      return { x: index, y: item.totalcpu }
   }).slice(-20);

   const options = {
      animationEnabled: true,
      exportEnabled: true,
      theme: "light2", // "light1", "dark1", "dark2"
      title: {
         text: "Uso de CPU"
      },
      axisY: {
         title: "CPU",
         suffix: "%"
      },
      axisX: {
         title: "TIMESTAMP",
         prefix: "T",
         interval: 2
      },
      data: [{
         type: "line",
         toolTipContent: "T {x}: {y} %",
         dataPoints: data2
      }]
   }
  
   return (
      <>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Monitor CPU
            </Typography>
            </center>
         </Paper>

         <br/>
         <br/>
         <div className='center' >
           
            <CanvasJSChart options={options} />
         </div>


         <br/>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Utilización
            </Typography>
            </center>
         </Paper>

         <br/>
         <center>
         <Typography variant="h5" color="inherit" component="div">
            {AllGenerales[AllGenerales.length-1].totalcpu} %
         </Typography>
         </center>

      </>
   );
}

export default Monitorcpu ;
