import { ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";

/**
 * v0 by Vercel.
 * @see https://v0.dev/t/Dp7EuNDak2F
 * Documentation: https://v0.dev/docs#integrating-generated-code-into-your-nextjs-app
 */
export default function ConsoleWidget() {

  const commandHistory = [
    {
      command: "cmd /c dir",
      response: "test.bat\npayload.exe",
    },
    {
      command: "cmd /c dir",
      response: "test.bat\npayload.exe",
    },
    {
      command: "cmd /c dir",
      response: "test.bat\npayload.exe",
    },
    {
      command: "cmd /c dir",
      response: "test.bat\npayload.exe",
    },
  ]

  const [expandedIndexes, setExpandedIndexes] = useState<any>([]);

  const toggleResponse = (index: any) => {
    setExpandedIndexes((prev: any) =>
      prev.includes(index)
        ? prev.filter((i: any) => i !== index) // collapse
        : [...prev, index] // expand
    );
  };

  return (
    // flex flex-col h-[100%]
    <div className="flex flex-col h-[100%] w-[100%] bg-[#0F172A] rounded-lg">
      <div className="flex flex-col h-[100%] overflow-y-auto p-1 gap-4 whitespace-pre-line">
        {commandHistory.map((command, index) => {
            return(
              <div className="space-y-2">
                <div onClick={() => toggleResponse(index)} className="flex items-center gap-2">
                  {expandedIndexes.includes(index) ?  <ChevronDown /> : <ChevronRight />}
                  <span className="text-[#9cdcfe]">voldemort@theagency:</span>
                  <span className="text-[#ce9178]">~$</span>
                  <span>{command.command}</span>
                  {/* <span className="animate-blink text-[#d4d4d4]">_</span> */}
                </div>

                {expandedIndexes.includes(index) && (
                  <div>
                    <span>{command.response}</span>
                  </div>
                )}
              </div>
            )
          })
        }
        
      </div>


      <div className="sticky bottom-0 bg-[#0F172A] p-4 flex items-center gap-2 rounded-lg">
        <span className="text-[#9cdcfe]">voldemort@theagency:</span>
        <span className="text-[#ce9178]">~$</span>
        <input
          type="text"
          className="bg-[#0F172A] outline-none flex-1 bottom-3 rounded-lg"
          placeholder="Enter command"
        />
      </div>
    </div>
  );
}
