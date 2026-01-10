# ArloCode Architecture Diagram

## System Overview

```mermaid
graph TB
    subgraph "Entry Point"
        Main[main.go]
        CMD[cmd/root.go]
    end

    subgraph "Terminal UI Layer"
        TUI[internal/tui/app]
        Layers[internal/tui/layers<br/>Z-axis Compositing]
        Notifs[internal/tui/notifications<br/>Animated Toasts]
        Conv[internal/tui/app/conversation<br/>Chat Manager]
        State[internal/tui/state<br/>Global State]
        Themes[internal/tui/themes]
    end

    subgraph "Agent Layer"
        CodingAgent[internal/coding_agent<br/>Pre-configured Agent]
    end

    subgraph "Butler Framework - Core"
        Agent[internal/butler/agent<br/>Orchestrator]
        Memory[internal/butler/memory<br/>Conversation History]
        Tools[internal/butler/tools<br/>Reflection-based Functions]
        EventHooks[internal/butler/types<br/>Event Hooks]
    end

    subgraph "Butler Framework - LLM Interface"
        LLMInterface[internal/butler/llm<br/>Provider Interface]
        OpenRouterLLM[internal/butler/llm/openrouter]
        GeminiLLM[internal/butler/llm/gemini]
        OpenAILLM[internal/butler/llm/openai]
    end

    subgraph "Provider Implementations"
        OpenRouterProvider[internal/butler/providers/openrouter]
        GeminiProvider[internal/butler/providers/gemini]
        OpenAIProvider[internal/butler/providers/openai]
        ProviderTypes[internal/butler/providers/providers.go<br/>ProviderResponse]
    end

    subgraph "External APIs"
        OpenRouterAPI[OpenRouter API<br/>Multiple Models]
        GeminiAPI[Google Gemini API]
        OpenAIAPI[OpenAI API]
    end

    subgraph "External Libraries"
        BubbleTea[Bubble Tea<br/>TUI Framework]
        Lipgloss[Lipgloss<br/>Styling]
        Harmonica[Harmonica<br/>Spring Physics]
        Cobra[Cobra<br/>CLI]
    end

    subgraph "Standard Toolset"
        ReadFile[read_file]
        ReadFolder[read_folder]
        ListContents[list_folder_contents]
        SearchCode[search_code]
        ApplyEdit[apply_edit]
        FetchURL[fetch_url_as_markdown]
        MakeFile[make_file]
        RunCommand[run_command]
    end

    %% Entry Flow
    Main --> CMD
    CMD --> TUI
    CMD --> CodingAgent
    CMD --> State

    %% TUI Internal
    TUI --> Layers
    TUI --> Notifs
    TUI --> Conv
    TUI --> Themes
    TUI --> State
    
    %% TUI to Agent
    TUI --> CodingAgent
    CodingAgent --> Agent

    %% Agent Core Relationships
    Agent --> Memory
    Agent --> Tools
    Agent --> LLMInterface
    Agent --> EventHooks
    EventHooks -.->|callbacks| TUI

    %% LLM Interface to Implementations
    LLMInterface --> OpenRouterLLM
    LLMInterface --> GeminiLLM
    LLMInterface --> OpenAILLM

    %% LLM to Providers
    OpenRouterLLM --> OpenRouterProvider
    GeminiLLM --> GeminiProvider
    OpenAILLM --> OpenAIProvider

    %% Providers to APIs
    OpenRouterProvider --> OpenRouterAPI
    GeminiProvider --> GeminiAPI
    OpenAIProvider --> OpenAIAPI

    %% Provider Response Flow
    OpenRouterProvider --> ProviderTypes
    GeminiProvider --> ProviderTypes
    OpenAIProvider --> ProviderTypes
    ProviderTypes --> LLMInterface

    %% Tools to Standard Toolset
    Tools --> ReadFile
    Tools --> ReadFolder
    Tools --> ListContents
    Tools --> SearchCode
    Tools --> ApplyEdit
    Tools --> FetchURL
    Tools --> MakeFile
    Tools --> RunCommand

    %% External Library Usage
    TUI --> BubbleTea
    TUI --> Lipgloss
    Notifs --> Harmonica
    CMD --> Cobra

    %% Styling
    classDef entry fill:#e78284,stroke:#c24a4f,stroke-width:3px,color:#fff
    classDef tui fill:#8caaee,stroke:#4673c4,stroke-width:2px,color:#fff
    classDef agent fill:#a6d189,stroke:#74a95d,stroke-width:2px,color:#000
    classDef butler fill:#ef9f76,stroke:#c97849,stroke-width:2px,color:#000
    classDef llm fill:#ca9ee6,stroke:#a275c4,stroke-width:2px,color:#fff
    classDef provider fill:#f4b8e4,stroke:#c488b4,stroke-width:2px,color:#000
    classDef external fill:#81c8be,stroke:#5b9c91,stroke-width:2px,color:#000
    classDef tools fill:#e5c890,stroke:#b59860,stroke-width:2px,color:#000
    classDef lib fill:#babbf1,stroke:#8a8bc1,stroke-width:2px,color:#000

    class Main,CMD entry
    class TUI,Layers,Notifs,Conv,State,Themes tui
    class CodingAgent agent
    class Agent,Memory,Tools,EventHooks butler
    class LLMInterface,OpenRouterLLM,GeminiLLM,OpenAILLM llm
    class OpenRouterProvider,GeminiProvider,OpenAIProvider,ProviderTypes provider
    class OpenRouterAPI,GeminiAPI,OpenAIAPI external
    class ReadFile,ReadFolder,ListContents,SearchCode,ApplyEdit,FetchURL,MakeFile,RunCommand tools
    class BubbleTea,Lipgloss,Harmonica,Cobra lib
```

## Data Flow Diagram

```mermaid
sequenceDiagram
    participant User
    participant TUI as TUI Layer
    participant Agent as Butler Agent
    participant Memory as Memory System
    participant LLM as LLM Interface
    participant Provider as Provider<br/>(OpenRouter/Gemini/OpenAI)
    participant Tools as Tool Executor

    User->>TUI: Types prompt
    TUI->>Agent: Run(prompt)
    Agent->>Memory: Add user message
    
    loop Iteration Loop (max 10-100)
        Agent->>LLM: Stream(memory, tools, hooks)
        LLM->>Provider: Convert & send request
        Provider-->>LLM: Stream response
        
        alt Thinking Stream
            LLM-->>TUI: OnThinkingChunk callback
            TUI-->>User: Display reasoning
        end
        
        alt Text Stream
            LLM-->>TUI: OnTextChunk callback
            TUI-->>User: Display response
        end
        
        alt Tool Calls
            LLM->>Agent: Return tool calls
            TUI-->>User: OnToolCall notification
            
            loop For each tool
                Agent->>Tools: Execute via reflection
                Tools-->>Agent: Return result
                Agent->>Memory: Add tool result
            end
            
            Note over Agent: Continue iteration
        end
        
        alt No More Tools
            LLM->>Agent: Final response
            Agent->>Memory: Add model response
            Agent-->>TUI: Complete
        end
    end
    
    TUI-->>User: Display final result
```

## Component Layer Architecture

```mermaid
graph LR
    subgraph "Layer 1: Presentation"
        A[Terminal UI<br/>Bubble Tea + Lipgloss]
        B[Notifications<br/>Harmonica Physics]
        C[Layers<br/>Z-compositing]
    end
    
    subgraph "Layer 2: Application"
        D[Coding Agent<br/>Configured Instance]
        E[State Management<br/>Global State]
    end
    
    subgraph "Layer 3: Business Logic"
        F[Butler Agent<br/>Orchestrator]
        G[Memory<br/>Conversation]
        H[Event Hooks<br/>Callbacks]
    end
    
    subgraph "Layer 4: Integration"
        I[LLM Interface<br/>Provider Abstraction]
        J[Tool System<br/>Reflection-based]
    end
    
    subgraph "Layer 5: External"
        K[AI Providers<br/>OpenRouter/Gemini/OpenAI]
        L[File System<br/>Tool Operations]
    end
    
    A --> D
    B --> D
    C --> A
    D --> F
    E --> D
    F --> G
    F --> H
    F --> I
    F --> J
    H -.->|callbacks| A
    I --> K
    J --> L
    
    classDef layer1 fill:#e78284,stroke:#c24a4f,stroke-width:2px,color:#fff
    classDef layer2 fill:#ef9f76,stroke:#c97849,stroke-width:2px,color:#000
    classDef layer3 fill:#a6d189,stroke:#74a95d,stroke-width:2px,color:#000
    classDef layer4 fill:#8caaee,stroke:#4673c4,stroke-width:2px,color:#fff
    classDef layer5 fill:#ca9ee6,stroke:#a275c4,stroke-width:2px,color:#fff
    
    class A,B,C layer1
    class D,E layer2
    class F,G,H layer3
    class I,J layer4
    class K,L layer5
```

## Tool Execution Flow

```mermaid
graph TB
    Start[AI Decides to Call Tool]
    
    Start --> Parse[Parse Tool Call<br/>JSON Arguments]
    Parse --> Lookup[Lookup Tool in Registry]
    Lookup --> Reflect[Use Reflection<br/>Get Handler Function]
    Reflect --> Marshal[Marshal Arguments<br/>to Struct]
    Marshal --> Invoke[Invoke Function<br/>with reflect.Value.Call]
    Invoke --> Execute[Execute Tool Logic]
    
    Execute --> FS{Tool Type?}
    
    FS -->|read_file| RF[Read from Filesystem]
    FS -->|search_code| SC[Search Files]
    FS -->|apply_edit| AE[Modify File]
    FS -->|make_file| MF[Create File]
    FS -->|run_command| RC[Execute Shell Command]
    FS -->|fetch_url| FU[HTTP Request]
    FS -->|list_folder| LF[Directory Listing]
    FS -->|read_folder| RFo[Read All Files]
    
    RF --> Return[Return Result String]
    SC --> Return
    AE --> Return
    MF --> Return
    RC --> Return
    FU --> Return
    LF --> Return
    RFo --> Return
    
    Return --> Memory[Add to Memory as<br/>Tool Result Entry]
    Memory --> Next[Continue Iteration]
    
    classDef decision fill:#f4b8e4,stroke:#c488b4,stroke-width:2px
    classDef process fill:#81c8be,stroke:#5b9c91,stroke-width:2px
    classDef tool fill:#e5c890,stroke:#b59860,stroke-width:2px
    classDef storage fill:#babbf1,stroke:#8a8bc1,stroke-width:2px
    
    class FS decision
    class Start,Parse,Lookup,Reflect,Marshal,Invoke,Return,Next process
    class RF,SC,AE,MF,RC,FU,LF,RFo tool
    class Memory,Execute storage
```

## Memory Structure

```mermaid
graph TB
    subgraph "Conversation Memory Array"
        M1[Entry 1<br/>Role: user<br/>Message: 'Read main.go']
        M2[Entry 2<br/>Role: model<br/>Message: ''<br/>ToolCalls: read_file]
        M3[Entry 3<br/>Role: tool<br/>Message: 'package main...'<br/>ToolCallID: call_123]
        M4[Entry 4<br/>Role: model<br/>Message: 'This file contains...']
        M5[Entry 5<br/>Role: user<br/>Message: 'Create a README']
        M6[Entry 6<br/>Role: model<br/>Message: ''<br/>ToolCalls: make_file]
        M7[Entry 7<br/>Role: tool<br/>Message: 'File created'<br/>ToolCallID: call_456]
        M8[Entry 8<br/>Role: model<br/>Message: 'Created README.md']
    end
    
    M1 --> M2 --> M3 --> M4 --> M5 --> M6 --> M7 --> M8
    
    M8 -.->|Next user input| M9[Entry 9<br/>Role: user]
    
    classDef user fill:#8caaee,stroke:#4673c4,stroke-width:2px,color:#fff
    classDef model fill:#a6d189,stroke:#74a95d,stroke-width:2px,color:#000
    classDef tool fill:#ef9f76,stroke:#c97849,stroke-width:2px,color:#000
    classDef future fill:#babbf1,stroke:#8a8bc1,stroke-width:2px,color:#000,stroke-dasharray: 5 5
    
    class M1,M5 user
    class M2,M4,M6,M8 model
    class M3,M7 tool
    class M9 future
```

---

## Legend

- **Red/Pink**: Entry points and main execution
- **Blue**: UI/TUI components
- **Green**: Agent layer
- **Orange**: Butler core framework
- **Purple**: LLM abstraction layer
- **Light Purple**: Provider implementations
- **Teal**: External services/APIs
- **Yellow**: Tool implementations
- **Light Blue**: External libraries

---

## Key Relationships

1. **TUI ← EventHooks → Agent**: Real-time streaming via callbacks
2. **Agent → Memory**: Every interaction stored for context
3. **Agent → LLM Interface**: Provider-agnostic communication
4. **LLM Interface → Providers**: Specific implementation adapters
5. **Agent → Tools**: Reflection-based dynamic execution
6. **Tools → Filesystem/Shell**: Actual operations on user's system

