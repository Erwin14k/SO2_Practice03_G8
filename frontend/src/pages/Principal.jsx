import React, {useState,useEffect} from 'react'

import Navbar from '../components/Navbar';
import Procesos from './Procesos';
import Monitorcpu from './Monitorcpu';
import Monitorram from './Monitorram';

const Principal=()=>{ 
   const [option, setOption] = useState(1);
   const [AllGenerales, setAllGenerales] = useState([]);
   const [AllProcesos, setAllProcesos] = useState([]);
   const [AllTipoProcesos, setAllTipoProcesos] = useState([]);

   const handleClick = (event) =>{
      event.preventDefault()
      setOption(event.target.textContent === "Procesos"? 1: event.target.textContent === "Monitor CPU"? 2: 3 )
   };

    
  useEffect(() => {
   console.log("===== START ENV =====")

      const fetchData = async () => {
         try{
            const response = await fetch(`http://${process.env.REACT_APP_PUERTO}:8080/tasks`);
            const data = await response.json();
            //console.log(data);
            
            setAllGenerales(data.AllGenerales);
            setAllProcesos(data.AllTipoProcesos );
            setAllTipoProcesos(data.AllProcesos);
            console.log("===== Get =====")
         }catch(error){
            console.log("===== Get Err =====")
            console.log(error)
            //alert(error)
         }
      }

      const interval = setInterval(fetchData, 3000);
      return () => clearInterval(interval);


   }, []);

   function conditionalRender(){
      switch (option) {
         case 1:
            return <Procesos AllProcesos={AllProcesos} AllTipoProcesos={AllTipoProcesos}  AllGenerales={AllGenerales}/>
         case 2:
            return <Monitorcpu AllGenerales={AllGenerales}/>
         case 3:
            return <Monitorram AllGenerales={AllGenerales} />
         default:
            return <Procesos/>
      }
   }

    return (
      <>
      <Navbar handleClick={handleClick}/>
      {conditionalRender()}
      </>
    );
}

export default Principal ;
