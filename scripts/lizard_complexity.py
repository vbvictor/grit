#!/usr/bin/env python3

"""
A script to run the Lizard code complexity analyzer and output results in CSV format.
Usage: 
  python lizard_complexity.py -l cpp -l c -t 8 path/to/code --output complexity.csv
"""

import argparse
import csv
import os
import sys
from typing import List, Dict, Any

# Import Lizard as a library
try:
    import lizard
except ImportError:
    print("Error: Lizard library not found. Install with 'pip install lizard'")
    sys.exit(1)


def parse_arguments() -> argparse.Namespace:
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(
        description="Generate code complexity CSV using Lizard"
    )
    parser.add_argument(
        "path", help="Path to analyze"
    )
    parser.add_argument(
        "-l", "--language", action="append", default=[],
        help="Programming language to include (can be specified multiple times)"
    )
    parser.add_argument(
        "-t", "--threads", type=int, default=1, 
        help="Number of threads to use for analysis"
    )
    parser.add_argument(
        "-m", "--modified", action="store_true", 
        help="Include only modified files"
    )
    parser.add_argument(
        "--output", default="complexity.csv", 
        help="Output CSV file path"
    )
    parser.add_argument(
        "--exclude", action="append", default=[], 
        help="File patterns to exclude (can be specified multiple times)"
    )
    parser.add_argument(
        "--verbose", action="store_true", 
        help="Enable verbose output"
    )
    
    return parser.parse_args()


def analyze_code(args: argparse.Namespace) -> List[Dict[str, Any]]:
    """Run Lizard analysis and return results."""
    lizard_args = []
    
    for lang in args.language:
        lizard_args.extend(["-l", lang])
    
    for pattern in args.exclude:
        lizard_args.extend(["--exclude", pattern])
    
    if args.threads > 1:
        lizard_args.extend(["-t", str(args.threads)])
    
    lizard_args.append(args.path)
    
    if args.verbose:
        print(f"Analyzing code in {args.path}...")
        print(f"Lizard arguments: {lizard_args}")
    
    analyzer = lizard.analyze(lizard_args)
    results = []
    
    for source_file in analyzer:
        file_path = os.path.relpath(source_file.filename, os.getcwd())
        
        for func in source_file.function_list:
            results.append({
                "file": file_path,
                "function": func.name,
                "length": func.length,
                "complexity": func.cyclomatic_complexity,
                "line": func.start_line,
            })
    
    return results


def write_csv(results: List[Dict[str, Any]], output_path: str, verbose: bool) -> None:
    """Write analysis results to CSV file."""
    with open(output_path, 'w', newline='') as csv_file:
        writer = csv.writer(csv_file)
        
        # Write data rows
        for result in results:
            writer.writerow([
                result["file"],
                result["function"],
                result["length"],
                result["complexity"],
                result["line"]
            ])
    
    if verbose:
        print(f"Wrote {len(results)} functions to {output_path}")


def main() -> None:
    """Main function."""
    args = parse_arguments()
    
    # Run analysis
    results = analyze_code(args)
    
    # Write results to CSV
    write_csv(results, args.output, args.verbose)
    
    if args.verbose:
        print(f"Analysis complete. Found {len(results)} functions.")


if __name__ == "__main__":
    main()
