function Show-Tree {
        param([string]$Path = ".")
            Get-ChildItem -Recurse $Path | ForEach-Object {
                        $depth = $_.FullName.Substring($Path.Length) -replace '[^\\]', '' | Measure-Object -Character
                                ('    ' * ($depth.Characters - 1)) + $_.Name
                                    }
}

Show-Tree

