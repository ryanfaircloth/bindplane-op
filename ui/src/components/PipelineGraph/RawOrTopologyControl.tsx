import { ToggleButton, ToggleButtonGroup } from "@mui/material";

interface RawOrTopologyControlProps {
  rawOrTopology: "raw" | "topology";
  setTopologyOrRaw: React.Dispatch<React.SetStateAction<"topology" | "raw">>;
}
export const RawOrTopologyControl: React.FC<RawOrTopologyControlProps> = ({
  rawOrTopology,
  setTopologyOrRaw,
}) => {
  
  const toggleRawOrTopology = (
    _: React.MouseEvent<HTMLElement>,
    newValue: "topology" | "raw",
    ) => {
      setTopologyOrRaw(newValue);
  };
  return (
    <ToggleButtonGroup
      sx={{
        display: 'flex',
        justifyContent: 'center',
      }}
      color="primary"
      value={rawOrTopology}
      exclusive
      onChange={toggleRawOrTopology}
    >
      <ToggleButton 
        value="topology" 
        sx={{
          display: "flex",
          justifyContent: "center",
          paddingLeft: 15,
          paddingRight: 15,
          textTransform: "none",
        }}
      >                    
        Topology
      </ToggleButton>
      <ToggleButton 
        value="raw" 
        sx={{
          display: "flex",
          justifyContent: "center",
          paddingLeft: 15,
          paddingRight: 15,
          textTransform: "none",
        }}
      >                    
        Raw
      </ToggleButton>                  
    </ToggleButtonGroup>
  )
};
      
      